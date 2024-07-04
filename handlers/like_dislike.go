// handlers/like_dislike.go

package handlers

import (
    "log"
    "net/http"
    "strconv"

    "DytForum/database"
)

// LikeThread handles HTTP requests to like or dislike a thread.
func LikeThread(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session-name")
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "You must be logged in to like or dislike a thread", http.StatusUnauthorized)
        return
    }

    userID := session.Values["user_id"].(int)
    threadID := r.FormValue("thread_id")
    likeType := r.FormValue("like_type")

    if threadID == "" || likeType == "" {
        http.Error(w, "Thread ID and like/dislike type are required", http.StatusBadRequest)
        return
    }

    likeTypeInt, err := strconv.Atoi(likeType)
    if err != nil {
        http.Error(w, "Invalid like/dislike type", http.StatusBadRequest)
        return
    }

    // Check if the user has already liked or disliked the thread
    var existingLike int
    err = database.DB.QueryRow("SELECT like_type FROM thread_likes WHERE thread_id = ? AND user_id = ?", threadID, userID).Scan(&existingLike)
    if err == nil {
        if existingLike == likeTypeInt {
            // User is trying to like or dislike the thread again with the same type, do nothing
            return
        } else {
            // User is changing their like/dislike type
            _, err := database.DB.Exec("UPDATE thread_likes SET like_type = ? WHERE thread_id = ? AND user_id = ?", likeTypeInt, threadID, userID)
            if err != nil {
                log.Printf("Failed to update like/dislike: %v", err)
                http.Error(w, "Failed to update like/dislike", http.StatusInternalServerError)
                return
            }
        }
    } else {
        // User hasn't liked or disliked the thread yet, insert new like/dislike
        _, err := database.DB.Exec("INSERT INTO thread_likes (thread_id, user_id, like_type) VALUES (?, ?, ?)", threadID, userID, likeTypeInt)
        if err != nil {
            log.Printf("Failed to create like/dislike: %v", err)
            http.Error(w, "Failed to create like/dislike", http.StatusInternalServerError)
            return
        }
    }

    // Update the thread's likes/dislikes count
    if likeTypeInt == 1 {
        _, err = database.DB.Exec("UPDATE threads SET likes = likes + 1 WHERE id = ?", threadID)
    } else {
        _, err = database.DB.Exec("UPDATE threads SET dislikes = dislikes + 1 WHERE id = ?", threadID)
    }
    if err != nil {
        log.Printf("Failed to update thread likes/dislikes: %v", err)
        http.Error(w, "Failed to update thread likes/dislikes", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

// LikeComment handles HTTP requests to like or dislike a comment.
func LikeComment(w http.ResponseWriter, r *http.Request) {
    session, _ := store.Get(r, "session-name")
    if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
        http.Error(w, "You must be logged in to like or dislike a comment", http.StatusUnauthorized)
        return
    }

    userID := session.Values["user_id"].(int)
    commentID := r.FormValue("comment_id")
    likeType := r.FormValue("like_type")

    if commentID == "" || likeType == "" {
        http.Error(w, "Comment ID and like/dislike type are required", http.StatusBadRequest)
        return
    }

    likeTypeInt, err := strconv.Atoi(likeType)
    if err != nil {
        http.Error(w, "Invalid like/dislike type", http.StatusBadRequest)
        return
    }

    // Check if the user has already liked or disliked the comment
    var existingLike int
    err = database.DB.QueryRow("SELECT like_type FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&existingLike)
    if err == nil {
        if existingLike == likeTypeInt {
            // User is trying to like or dislike the comment again with the same type, do nothing
            return
        } else {
            // User is changing their like/dislike type
            _, err := database.DB.Exec("UPDATE comment_likes SET like_type = ? WHERE comment_id = ? AND user_id = ?", likeTypeInt, commentID, userID)
            if err != nil {
                log.Printf("Failed to update like/dislike: %v", err)
                http.Error(w, "Failed to update like/dislike", http.StatusInternalServerError)
                return
            }
        }
    } else {
        // User hasn't liked or disliked the comment yet, insert new like/dislike
        _, err := database.DB.Exec("INSERT INTO comment_likes (comment_id, user_id, like_type) VALUES (?, ?, ?)", commentID, userID, likeTypeInt)
        if err != nil {
            log.Printf("Failed to create like/dislike: %v", err)
            http.Error(w, "Failed to create like/dislike", http.StatusInternalServerError)
            return
        }
    }

    // Update the comment's likes/dislikes count
    if likeTypeInt == 1 {
        _, err = database.DB.Exec("UPDATE comments SET likes = likes + 1 WHERE id = ?", commentID)
    } else {
        _, err = database.DB.Exec("UPDATE comments SET dislikes = dislikes + 1 WHERE id = ?", commentID)
    }
    if err != nil {
        log.Printf("Failed to update comment likes/dislikes: %v", err)
        http.Error(w, "Failed to update comment likes/dislikes", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
