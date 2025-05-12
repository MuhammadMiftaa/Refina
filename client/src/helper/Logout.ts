import Cookies from "js-cookie";

export const HandleLogout = () => {
  const confirmation = confirm("Are you sure you want to logout?");
  if (!confirmation) return;

  Cookies.remove("token");
  window.location.href = "/login";
};
