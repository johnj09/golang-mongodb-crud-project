import { useEffect, useContext } from "react";
import { Link } from "react-router-dom";
import { UserContext } from "../util/UserContext";

function Header() {
  const { setUserInfo, userInfo } = useContext(UserContext);

  useEffect(() => {
    fetch(process.env.REACT_APP_API_BASE_URL + "/api/profile", {
      credentials: "include",
    }).then((res) => {
      if (res.ok) {
        res.json().then((data) => {
          setUserInfo(data.data.data);
        });
      } else {
        setUserInfo(null);
      }
    });
  }, []);

  const handleLogout = () => {
    fetch(process.env.REACT_APP_API_BASE_URL + "/api/logout", {
      credentials: "include",
      method: "POST",
    });
    setUserInfo(null);
  };

  const username = userInfo?.username;

  return (
    <header>
      <Link to="/" className="logo">
        MyBlog
      </Link>
      <nav>
        {username && (
          <>
            <Link to="/create">Create Post</Link>
            <a className="logout" onClick={handleLogout}>
              Logout
            </a>
          </>
        )}
        {!username && (
          <>
            <Link to="/login">Login</Link>
            <Link to="/register">Register</Link>
          </>
        )}
      </nav>
    </header>
  );
}

export default Header;
