import { LiaComment, LiaHeart } from "react-icons/lia";
import { Link } from "react-router-dom";

function Post({ id, title, content, user_id, created_at, num_of_comments }) {
  const date = new Date(created_at * 1000); // Convert seconds to milliseconds
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0"); // Months are zero-based
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");
  const seconds = String(date.getSeconds()).padStart(2, "0");
  const timeCreated = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;

  return (
    <div className="post">
      <Link to={`/post/${id}`}>
        <h2 className="post-title">{title}</h2>
      </Link>
      <p className="post-info">
        <a href="" className="author">
          User {user_id}
        </a>
        <time datetime="">{timeCreated}</time>
        {/* <span className="likes">
          <LiaHeart className="heart-icon" />
          <p>10</p>
        </span> */}
        <span className="comments">
          <LiaComment className="comment-icon" />
          <p>{num_of_comments}</p>
        </span>
      </p>
      <p className="post-summary">
        {content.length > 100 && content.substring(0, 100) + "..."}
        {content.length <= 100 && content}
      </p>
    </div>
  );
}

export default Post;
