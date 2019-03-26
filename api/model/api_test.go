package model

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var tempvid string

func clearTables() {
	DbConn.Exec("truncate users")
	DbConn.Exec("truncate video_info")
	DbConn.Exec("truncate comments")
	DbConn.Exec("truncate sessions")
}

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}

func TestUserWorkFlow(t *testing.T) {
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUser)
	t.Run("Del", testDeleteUser)
	t.Run("Reget", testRegetUser)
}

func testAddUser(t *testing.T) {
	if err := AddUserCredential("zhangsan", "123"); err != nil {
		t.Errorf("Error of AddUser:%v", err)
	}
}

func testGetUser(t *testing.T) {
	if pwd, err := GetUserCredential("zhangsan"); pwd != "123" || err != nil {
		t.Errorf("Error of GetUser:%s", err)
	}
}

func testDeleteUser(t *testing.T) {
	if err := DeleteUser("zhangsan", "123"); err != nil {
		t.Errorf("Error of DeleteUser:%v", err)
	}
}

func testRegetUser(t *testing.T) {
	if pwd, err := GetUserCredential("zhangsan"); err != nil {
		t.Errorf("Error of RegetUser:%v", err)
	} else if pwd != "" {
		t.Errorf("Deleting user test failed")
	}
}

func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("PrepareUser", testAddUser)
	t.Run("AddVideo", testAddVideoInfo)
	t.Run("GetVideo", testGetVideoInfo)
	t.Run("DelVideo", testDeleteVideoInfo)
	t.Run("RegetVideo", testRegetVideoInfo)
}

func testAddVideoInfo(t *testing.T) {
	if vi, err := AddNewVideo(1, "my-video"); err != nil {
		t.Errorf("Error of AddVideoInfo:%v", err)
	} else {
		tempvid = vi.Id
	}
}

func testGetVideoInfo(t *testing.T) {
	if _, err := GetVideoInfo(tempvid); err != nil {
		t.Errorf("Error of GetVideoInfo:%v", err)
	}
}

func testDeleteVideoInfo(t *testing.T) {
	if err := DeleteVideoInfo(tempvid); err != nil {
		t.Errorf("Error of DeleteVideoInfo: %v", err)
	}
}

func testRegetVideoInfo(t *testing.T) {
	if vi, err := GetVideoInfo(tempvid); err != nil || vi != nil {
		t.Errorf("Error of RegetVideoInfo:%v", err)
	}
}

func TestComments(t *testing.T) {
	clearTables()
	t.Run("AddUser", testAddUser)
	t.Run("AddComments", testAddComments)
	t.Run("ListComments", testListComments)
}

func testAddComments(t *testing.T) {
	vid := "12345"
	aid := 1
	content := "I like this video"

	if err := AddNewComments(vid, aid, content); err != nil {
		t.Errorf("Error of AddComments:%v", err)
	}
}

func testListComments(t *testing.T) {
	vid := "12345"
	from := 1514764800
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000, 10))

	if res, err := ListComments(vid, from, to); err != nil {
		t.Errorf("Error of ListComments:%v", err)
	} else {
		for i, ele := range res {
			fmt.Printf("comment:%d,%v\n", i, ele)
		}
	}
}
