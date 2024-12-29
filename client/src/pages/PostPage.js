import { useContext, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { UserContext } from "../util/UserContext";
import { UnixToDate } from "../util/util";
import Comment from "../components/Comment";

function PostPage() {
  const [postInfo, setPostInfo] = useState(null);
  const [comments, setComments] = useState([]);
  const { id } = useParams();
  const { userInfo } = useContext(UserContext);
  useEffect(() => {
    fetch(process.env.REACT_APP_API_BASE_URL + `/api/post/${id}`, {
      method: "GET",
    })
      .then((res) => {
        res.json().then((data) => {
          console.log(data.data.data);
          setPostInfo(data.data.data);
        });
      })
      .catch((err) => console.error(err));
    fetch(process.env.REACT_APP_API_BASE_URL + `/api/comment/${id}`, {
      method: "GET",
    })
      .then((res) => {
        res.json().then((data) => {
          console.log(data.data.data);
          setComments(data.data.data);
        });
      })
      .catch((err) => console.error(err));
  }, []);

  if (!postInfo) return "";

  return (
    <div className="post-page-container">
      <h1 className="post-page-title">{postInfo.title}</h1>
      <div className="post-page-meta">
        <span className="post-page-author">By {postInfo.user_id}</span>
        <span className="post-page-time">
          Created at {UnixToDate(postInfo.created_at)}
        </span>
        <span className="post-page-comments">
          {postInfo.num_of_comments} Comments
        </span>
      </div>
      <p className="post-page-content">{postInfo.content}</p>
      <div className="post-page-comments-section">
        <h2 className="comments-section-title">Comments</h2>
        {comments ? (
          <ul className="comments-list">
            {comments.map((comment, index) => (
              <Comment {...comment} />
            ))}
          </ul>
        ) : (
          <p className="no-comments">No comments yet.</p>
        )}
      </div>
    </div>
  );
}

export default PostPage;
