import { useEffect, useState } from "react";
import Post from "../components/Post";

function HomePage() {
  const [posts, setPosts] = useState([]);

  useEffect(() => {
    fetch(process.env.REACT_APP_API_BASE_URL + "/api/posts/0", {
      method: "GET",
      credentials: "include",
    }).then((res) => {
      res.json().then((data) => {
        setPosts(data.data.data);
        console.log(data.data.data);
      });
    });
  }, []);

  return <>{posts.length > 0 && posts.map((post) => <Post {...post} />)}</>;
}
export default HomePage;
