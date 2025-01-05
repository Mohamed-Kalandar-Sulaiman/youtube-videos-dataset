Module Name : Mohamed-Kalandar-Sulaiman/youtube-videos-dataset

# Steps to run the project
1. To run it locally clone the project 
2. Execute the command : docker-compose up --build

Curl command to access the api
curl --location 'localhost:3000/api/v1/videos?pageSize=2&PublishedBefore=2025-01-05T00%3A00%3A00Z&channelName=ku&PublishedAfter=2025-01-04T00%3A00%3A00Z&pageKey=2025-01-03T00%3A15%3A01Z&title=null'