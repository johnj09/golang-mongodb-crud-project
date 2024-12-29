import axios from "axios";

const API = axios.create({
  baseURL: process.env.REACT_APP_API_BASE_URL, // Replace with your backend URL
  withCredentials: true, // To include cookies
});

// Intercept requests to include the token in headers
// API.interceptors.request.use((req) => {
//   const token = localStorage.getItem("token");
//   if (token) {
//     req.headers.Authorization = `Bearer ${token}`;
//   }
//   return req;
// });

// Interceptors for debugging
API.interceptors.request.use(
  (config) => {
    console.log("Making request to:", config.baseURL + config.url);
    return config;
  },
  (error) => {
    console.error("Request error:", error);
    return Promise.reject(error);
  }
);

API.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    console.error("Response error:", error.response || error.message);
    return Promise.reject(error);
  }
);

// Post APIs
export const fetchPosts = (page = 0) => API.get(`/posts/${page}`);
export const fetchPost = (pid) => API.get(`/post/${pid}`);
export const createPost = (data) => API.post(`/post`, data);
export const updatePost = (pid, data) => API.put(`/post/${pid}`, data);
export const deletePost = (pid) => API.delete(`/post/${pid}`);

// Comment APIs
export const createComment = (pid, data) => API.post(`/comment/${pid}`, data);
export const fetchComments = (pid) => API.get(`/comment/${pid}`);
export const updateComment = (cid) => API.patch(`/comment/${cid}`);
export const deleteComment = (cid) => API.delete(`/comment/${cid}`);

// Auth APIs
export const login = (loginReq) => API.post(`/login`, loginReq);
export const register = (registerReq) => API.post(`/api/register`, registerReq);

export default API;
