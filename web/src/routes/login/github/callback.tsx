import { createFileRoute } from "@tanstack/react-router";
import axios from "axios";
import { useEffect } from "react";
import { API_URL } from "../../../hooks/react-query/constants";

export const Route = createFileRoute("/login/github/callback")({
  component: Callback,
});

const getToken = async (code: string) => {
  const tok = await axios.get(
    `${API_URL}/auth/login/github/callback?code=${code}`
  );
  console.log(tok.data);
};

function Callback() {
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get("code");
    console.log(`|${code}|`);
    if (code) {
      getToken(code).then().catch();
    }
  }, []);
  return <></>;
}
