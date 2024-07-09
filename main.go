package main

import (
	"os"

	"fmt"
	"log"
	"net/http"
	"html/template"

	"encoding/json"
	"io/ioutil"

	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

var siteTitle = "LOGOS"

var wikiDataDir = "./wikidata/src/"

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	r := gin.Default()

	users := make(map[string]string)
	// Read the JSON file into a byte slice
	data, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	// Parse the byte slice into the users map
	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Fatalf("unable to parse user data: %v", err)
	}

	// Setup session middleware
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("tmpl/*")

	r.Static("/static", "./static")
	r.Static("/wikidata/dst", "./w")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"siteTitle": siteTitle,
		})
	})

	r.POST("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		username := c.PostForm("username")
		password := c.PostForm("password")

		expectedPassword, ok := users[username]
		if !ok || expectedPassword != password {
			c.JSON(401, gin.H{"status": "unauthorized"})
			return
		}

		session.Set("user", username)
		session.Save()

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"siteTitle": siteTitle,
			"port":      PORT,
			"addr":      "localhost",
			"message":   "Logged in successfully",
		})
	})

	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			// c.JSON(401, gin.H{"status": "unauthorized"}) // Not logged in
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"siteTitle": siteTitle,
				"port":      PORT,
				"addr":      "localhost",
			})
			return
		}

		// Delete the user session (logout)
		session.Delete("user")
		session.Save()

		// c.JSON(200, gin.H{"status": "logged out"})
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"siteTitle": siteTitle,
			"port":      PORT,
			"addr":      "localhost",
			"alert":     "Logged out successfully",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"siteTitle": siteTitle,
			"port":      PORT,
			"addr":      "localhost",
		})
	})

	r.POST("/submit/:title", authRequired(), func(c *gin.Context) {
    // Get the title and text from the form
    title := c.Param("title")
    textContent := c.PostForm("text")

    if !secureTitle(title) {
			c.String(400, "Title must be alphanumeric")
			return
    }

    if err := saveNoweb(textContent, wikiDataDir, title); err != nil {
			c.String(500, fmt.Sprintf("Failed to save file: %v", err))
			return
    }

    // Respond with a success message
    c.String(200, "File saved: %s.nw", title)
	})
	
	r.GET("/submit/:title", authRequired(), func(c *gin.Context) {
		title := c.Param("title")
		c.HTML(http.StatusOK, "newpage.tmpl", gin.H{
			"title": title,
		})
	})

	r.GET("/edit/:title", authRequired(), func(c *gin.Context) {
		title := c.Param("title")

    if !secureTitle(title) {
        c.String(400, "Title must be alphanumeric")
        return
    }

		filePath := filepath.Join(wikiDataDir, title+".nw")
		
		content, err := ioutil.ReadFile(filePath)
    if err != nil {
        c.String(500, fmt.Sprintf("Failed to read file: %v", err))
        return
    }
		
		c.HTML(http.StatusOK, "editpage.tmpl", gin.H{
			"title": title,
			"content": template.HTML(content),
		})
	})

	r.Run(":" + PORT)
}

func secureTitle(s string) bool {
    re := regexp.MustCompile("^[a-zA-Z0-9]+$")
    return re.MatchString(s)
}

func saveNoweb(textContent, dir, filename string) error {
    if err := os.MkdirAll(dir, os.ModePerm); err != nil {
        return err
    }

    filePath := filepath.Join(dir, filename+".nw")

    if err := ioutil.WriteFile(filePath, []byte(textContent), 0644); err != nil {
        return err
    }

    return nil
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.JSON(401, gin.H{"status": "unauthorized"})
			c.Abort()
			return
		}
		// If user is found, pass to the next middleware
		c.Next()
	}
}
