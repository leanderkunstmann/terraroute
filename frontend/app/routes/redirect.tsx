import { useEffect } from "react";
import { useNavigate } from "react-router";

export default function Redirect() {
  const navigate = useNavigate();

  useEffect(() => {
    navigate("/");
  }, [navigate]);

  return null;
}
