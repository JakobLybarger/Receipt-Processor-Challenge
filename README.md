# Receipt Processor

### Installation
1. Clone the repository
```bash
git clone https://github.com/JakobLybarger/Receipt-Processor-Challenge.git
cd Receipt-Processor-Challenge
```
2. Install dependencies:
```bash
go mod tidy
```
3. Run the app locally:
```bash
go run main.go


### Running with Docker
1. Build the docker image:
```bash
docker build -t receipt-processing
```
2. Run the container:
```bash
docker run -p 8080:8080 receipt-processing
```
