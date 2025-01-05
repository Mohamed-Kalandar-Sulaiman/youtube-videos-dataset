# ***Basic Requirements:***

    1. Fetch youtube videos data in backgorund
        Data :
            Search Query : LeetCode
            Fields:
                videoTitle
                videoDescription
                videoPublished_date
                videoThumbnail_url
                ... Additional fields
                videoLikes
                channelName
                videoDuration


    2. Provide following API to access the fetched data, [Decided to make a single endpoint Since both are search apis , Default behaviour is Point 2]
        1. GET /videos
        
        Query Params : 
            1. pageSize
            2. PageKey
            3. title
            4. channel_name
            5. published_date
            6. 

        Default Returns > Paginated set of video data sorted by published_date ORDER BY DESC
       
        Example:
        {
            "data" : {
                [List of Videos]
            },
            "pagination":{
                "nextPageKey:"<....>",
                "lastPageKey": "<....>
            }
        }

        - Use cursor based pagination

    
    
# ***Additional Points***
1. Deploy the app
2. Fetch enough data points
    - Maybe collect data locally
    - Deploy with sample data in AWS !
3. Postman collection
4. 


# Decisions
1. Use Fiber with Postgres
2. Postgres
    - Read replicas ?
    - Partition rows based on published_date
    