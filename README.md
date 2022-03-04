# url-shortner with fiber and redis
To run the program make sure you are in the root directory of project and run this command : 

```
docker-compose up 
```
This project has two main url
| url     | functionality |
| ---      | ---       |
| localhost:3000/:url | Get method, if url exist in db redirect to main url|
| localhost:3000/api/va1 | Post method, user send url,expiry and short parameters with body        |
