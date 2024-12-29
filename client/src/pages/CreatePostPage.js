import React, { useState } from "react";
import { Navigate } from "react-router-dom";

const CreatePostPage = () => {
  // State to hold the title and content of the post
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [redirect, setRedirect] = useState(false);

  // Handle form submission
  const handleSubmit = async (e) => {
    e.preventDefault();

    // Simulate post submission logic
    const postData = {
      title,
      content,
    };

    await fetch(process.env.REACT_APP_API_BASE_URL + "/api/post", {
      method: "POST",
      body: JSON.stringify(postData),
      headers: { "Content-Type": "application/json" },
      credentials: "include",
    })
      .then((res) => {
        if (res.ok) {
          alert("Post Successfully Created.");
          setRedirect(true);
        } else {
          alert("There was an issue while creating the post.");
        }
      })
      .catch((err) => alert("There was an issue while creating the post."));
  };

  if (redirect) {
    return <Navigate to={"/"} />;
  }

  return (
    <div className="create-post-container">
      <form onSubmit={handleSubmit} className="create-post-form">
        <label htmlFor="post-title" className="create-post-label">
          Title
        </label>
        <input
          id="post-title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Enter the title"
          required
          className="create-post-input"
        />

        <label htmlFor="post-content" className="create-post-label">
          Content
        </label>
        <textarea
          id="post-content"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder="write your post content here..."
          rows="10"
          required
          className="create-post-textarea"
        />

        <button type="submit" className="create-post-button">
          Submit
        </button>
      </form>
    </div>
  );
};

export default CreatePostPage;
