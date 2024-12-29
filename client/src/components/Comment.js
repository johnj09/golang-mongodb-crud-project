function Comment({ id, user_id, created_at, content }) {
  return (
    <li key={id} className="comment-item">
      <span className="comment-author">{user_id}</span>: {content}
    </li>
  );
}

export default Comment;
