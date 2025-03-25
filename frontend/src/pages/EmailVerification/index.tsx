import { useEffect, useState } from "react";
import { MailCheck } from "lucide-react";

import PageLoader from "../../components/PageLoader";
import NotFound from "../NotFound";

import useAxios from "../../hooks/useAxios";
import {
  API_METHOD_PATCH,
  API_VERIFY_EMAIL,
  EMAIL_VERIFY_TITLE,
  TOKEN_KEY,
} from "../../const/const";
import { Response } from "../../types/response";
import useTitle from "../../hooks/useTitle";

export default function EmailVerification() {
  const [isSuccess, setIsSuccess] = useState<boolean>(false);

  useTitle(EMAIL_VERIFY_TITLE)

  const queryParams = new URLSearchParams(window.location.search);
  const token = queryParams.get(TOKEN_KEY);

  if (!token) {
    return <NotFound />;
  }

  const { error, isLoading, fetchData } = useAxios<Response<undefined>>(
    API_VERIFY_EMAIL,
    API_METHOD_PATCH
  );

  useEffect(() => {
    fetchData({ token: token }).then((res) => {
      if (res && res.message) {
        setIsSuccess(true);
      }
    });
  }, []);

  if (isLoading) {
    return <PageLoader />;
  }

  if (error) {
    return <NotFound />;
  }

  return (
    <div className="grow flex bg-primary-white">
      {isSuccess && (
        <div className="w-[90%] lg:w-[80%] mx-auto flex flex-col justify-center items-center gap-3">
          <MailCheck
            className="h-32 w-32 md:h-64 md:w-64 text-primary-green"
            strokeWidth={0.5}
          />
          <h1 className="text-center text-2xl md:text-4xl font-bold text-primary-black">
            Email Verified
          </h1>
          <p className="text-center md:text-xl text-primary-black">
            Your email address is now verified
          </p>
        </div>
      )}
    </div>
  );
}
