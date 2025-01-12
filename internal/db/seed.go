package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/ngvanthanggit/RicolaSocial/internal/store"
)

const numPosts = 10
const numUsers = 5
const numComments = 20

var titleSample []string = []string{
	"Boost Your Day!",
	"Stay Motivated!",
	"Morning Wins!",
	"Travel Goals 2025",
	"Embrace Minimalism",
	"Balance Your Life",
	"Start Journaling",
	"Quick Recipes",
	"Save Smarter",
	"Build Habits",
	"Tech Trends",
	"DIY Ideas",
	"Books to Read",
	"Think Positive",
	"Daily Self-Care",
	"Network Better",
	"Weekend Hacks",
	"Sleep Tips",
	"Quick Workouts",
	"Cool Gadgets",
}

var contentSample []string = []string{
	"Start your day with a smile and a cup of coffee.",
	"Did you know taking a 10-minute walk can boost your mood instantly?",
	"Transform your space with these three simple DIY tricks.",
	"Ready for a challenge? Try this five-minute workout today.",
	"What's your dream travel destination? Share in the comments.",
	"Organize your life with these three easy tips.",
	"Discover the latest tech gadgets of 2025.",
	"Five reasons why journaling can change your life.",
	"Cooking made easy: Three-ingredient recipes to try tonight.",
	"Start saving for your goals with these simple strategies.",
	"New book alert! Here are the top three must-reads this month.",
	"How do you unwind after a busy day? Let us know.",
	"Create lasting habits with this simple formula.",
	"Quick tip: Drink water first thing in the morning.",
	"Find joy in the little things. What made you smile today?",
	"Power up your day with these healthy snack ideas.",
	"Weekend vibes: What's your go-to activity for relaxation?",
	"Learn how to stay focused with this one-minute breathing exercise.",
	"Say goodbye to clutter with these minimalist tips.",
	"Discover how to build connections that last.",
}

var tagSample []string = []string{
	"#Motivation", "#LifeHacks", "#Fitness", "#Travel", "#TechTrends", "#DIYProjects", "#Wellness",
	"#FoodLovers", "#PersonalGrowth", "#BookRecommendations", "#MoneyTips", "#Relaxation", "#HealthTips",
	"#Minimalism", "#Networking", "#WeekendVibes", "#Productivity", "#Mindfulness", "#SelfCare", "#Happiness",
}

var commentSample []string = []string{
	"Great post!", "Love this!", "Amazing content!", "So inspiring!", "Thanks for sharing!",
	"This is so true!", "Couldn't agree more.", "Very helpful, thank you!", "Keep it up!",
	"This made my day.", "Wow, just wow!", "Interesting perspective.", "Totally relatable!",
	"This is gold!", "Great tips!", "I needed this today.", "Well said!", "You’re so right.",
	"Love your content!", "Such good vibes!", "This is exactly what I was looking for.",
	"Mind blown!", "Can’t wait to try this.", "Super useful, thanks!", "This is brilliant!",
	"Absolutely love this.", "So creative!", "Thanks for the advice.", "Such a good idea.",
	"Keep the posts coming!",
}

var usernameSample []string = []string{
	"alice", "bob", "charlie", "dave", "eve", "frank", "grace", "hannah", "ivan", "julia",
	"kate", "leo", "mia", "nate", "olivia", "peter", "quinn", "rose", "sam", "tina",
	"uma", "vick", "will", "xena", "yara", "zane", "abby", "ben", "chris", "diana",
	"ella", "finn", "gina", "harry", "iris", "jack", "kyle", "luna", "mike", "nina",
	"oscar", "paul", "queen", "ryan", "sara", "tom", "ursula", "violet", "wade", "zoe",
}

// handler for generating new data for users, posts, comments
func SeedHandler(store store.Storage) {
	ctx := context.Background()

	// generate users
	users := generateUsers()
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user:", err)
			return
		}
	}

	// generate posts
	posts := generatePosts(users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post:", err)
			return
		}
	}

	//generate comments
	comments := generateComments(users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return
		}
	}

	log.Println("Seeding completed!")
}

func generateUsers() []*store.User {
	users := make([]*store.User, numUsers)

	for i := 0; i < numUsers; i++ {
		username := usernameSample[rand.Intn(len(usernameSample))] + fmt.Sprintf("%d", i)
		users[i] = &store.User{
			Username: username,
			Email:    username + "@example.com",
			Password: "12345",
		}
	}

	return users
}

func generatePosts(users []*store.User) []*store.Post {
	posts := make([]*store.Post, numPosts)

	for i := 0; i < numPosts; i++ {
		posts[i] = &store.Post{
			Title:   titleSample[rand.Intn(len(titleSample))],
			Content: contentSample[rand.Intn(len(contentSample))],
			Tags: []string{
				tagSample[rand.Intn(len(tagSample))],
				tagSample[rand.Intn(len(tagSample))],
			},
			UserID: users[rand.Intn(len(users))].ID,
		}
	}

	return posts
}

func generateComments(users []*store.User, posts []*store.Post) []*store.Comment {
	comments := make([]*store.Comment, numComments)

	for i := 0; i < numComments; i++ {
		comments[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: commentSample[rand.Intn(len(commentSample))],
		}
	}

	return comments
}
