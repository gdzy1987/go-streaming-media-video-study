## API接口
- 测试

1.创建用户接口:CreateUser POST localhost:8000/user
```
- body
{
	"user_name":"lisi",
	"pwd":"123456"
}

- return
{
    "success": true,
    "session_id": "f1b8d11b-5332-4860-8bef-c00f9826ec86"
}
```
2.登陆接口:Login POST localhost:8000/user/lisi
```
- body
{
	"user_name":"lisi",
	"pwd":"123456"
}

- return 
{
    "success": true,
    "session_id": "25cbca33-bdfc-412a-8735-fdeaed6641eb"
}
```
3.获取用户ID:GetUserInfo GET localhost:8000/user/lisi
```
- headers
X-User-Name:f1b8d11b-5332-4860-8bef-c00f9826ec86

- return 
{
    "id": 5
}
```
4.上传视频:AddNewVideo POST localhost:8000/user/lisi/videos
```
- body
{
	"author_id":5,
	"name":"视频1"
}

- headers
X-User-Name:f1b8d11b-5332-4860-8bef-c00f9826ec86

- return 

{
    "id": "ea893867-9dfc-4464-ada9-7c940eb6e277",
    "author_id": 5,
    "name": "视频1",
    "display_ctime": "Mar 07 2019, 11:51:07"
}
```

5.获取用户上传的视频:ListAllVideos  GET localhost:8000/user/lisi/videos
```
- headers
X-User-Name:f1b8d11b-5332-4860-8bef-c00f9826ec86

- return 
{
    "videos": [
        {
            "id": "747fdaa8-34f6-4dee-92f5-e0a59d3ebfdf",
            "author_id": 5,
            "name": "视频1",
            "display_ctime": "Mar 07 2019, 11:57:26"
        },
        {
            "id": "ea893867-9dfc-4464-ada9-7c940eb6e277",
            "author_id": 5,
            "name": "视频1",
            "display_ctime": "Mar 07 2019, 11:51:07"
        }
    ]
}
```
6.删除用户视频:DeleteVideo DELETE  localhost:8000/user/lisi/videos/747fdaa8-34f6-4dee-92f5-e0a59d3ebfdf
```
- headers
X-User-Name:f1b8d11b-5332-4860-8bef-c00f9826ec86

- return

```

7.发布评论: PostComment POST localhost:8000/videos/1d307d4b-3e95-42a0-975f-75e8facccca6/comments
```
- body
{
	"author_id":5,
	"content":"视频1真好看"
}
- headers
X-User-Name:f1b8d11b-5332-4860-8bef-c00f9826ec86

- return
```
8.获取评论: ShowComments GET localhost:8000/videos/1d307d4b-3e95-42a0-975f-75e8facccca6/comments
```
- headers
X-User-Name:f1b8d11b-5332-4860-8bef-c00f9826ec86

- return
{
    "comments": [
        {
            "id": "7f892c2d-730b-4c18-9896-036341f38e6e",
            "video_id": "1d307d4b-3e95-42a0-975f-75e8facccca6",
            "author": "lisi",
            "content": "视频1真好看"
        },
        {
            "id": "bab52c6c-cece-41cb-82b9-c08a4f048f6e",
            "video_id": "1d307d4b-3e95-42a0-975f-75e8facccca6",
            "author": "lisi",
            "content": "视频1真好看"
        }
    ]
}
```