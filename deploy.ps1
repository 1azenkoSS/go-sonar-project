[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

Write-Host "STARTING CI/CD PIPELINE..." -ForegroundColor Cyan

#TESTING
Write-Host "1. Running Unit Tests..." -ForegroundColor Yellow
go test -coverprofile coverage.out ./...

if ($LASTEXITCODE -ne 0) {
    Write-Host "Tests Failed! Stopping pipeline." -ForegroundColor Red
    exit 1
}
Write-Host "Tests Passed!" -ForegroundColor Green

#QUALITY ANALYSIS (SONARQUBE)
Write-Host "2. Running SonarQube Analysis..." -ForegroundColor Yellow
$SONAR_TOKEN = "sqp_66dd5187a42de41868a87ba8a6ab71a710434c70" 

docker run --rm -v "${PWD}:/usr/src" --network=go-sonar-project_sonarnet -e SONAR_HOST_URL="http://sonarqube:9000" -e SONAR_TOKEN="$SONAR_TOKEN" sonarsource/sonar-scanner-cli

if ($LASTEXITCODE -ne 0) {
    Write-Host "SonarQube Scan Failed!" -ForegroundColor Red
    exit 1
}
Write-Host "Code Analysis Uploaded!" -ForegroundColor Green

#BUILD
Write-Host "3. Building Docker Image..." -ForegroundColor Yellow


$DOCKER_USER = "yuriilazenko" 
$IMAGE_NAME = "$DOCKER_USER/credit-calculator:latest"

docker build -t $IMAGE_NAME .

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build Failed!" -ForegroundColor Red
    exit 1
}
Write-Host "Docker Image Built Successfully!" -ForegroundColor Green


Write-Host "3.1. Pushing to Docker Hub..." -ForegroundColor Yellow

docker push $IMAGE_NAME

if ($LASTEXITCODE -ne 0) {
    Write-Host "Push Failed! (Did you run 'docker login'?)" -ForegroundColor Red

} else {
    Write-Host "Image Pushed to Docker Hub!" -ForegroundColor Green
}

#DEPLOY
Write-Host "4. Deploying Application..." -ForegroundColor Yellow

docker stop go-app-container 2>$null
docker rm go-app-container 2>$null

docker run -d --name go-app-container -p 8080:8080 $IMAGE_NAME

Write-Host "DEPLOYMENT COMPLETE!" -ForegroundColor Magenta
Write-Host "Your app is running at: http://localhost:8080" -ForegroundColor Cyan