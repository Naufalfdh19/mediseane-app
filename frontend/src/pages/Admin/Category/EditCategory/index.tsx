import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { createPortal } from "react-dom";

import Input from "../../../../components/Input";
import Button from "../../../../components/Button";
import ConfirmBox from "../../../../components/ConfirmBox";
import LoaderBox from "../../../../components/LoaderBox";
import ErrorBox from "../../../../components/ErrorBox";
import SuccessBox from "../../../../components/SuccessBox";
import TitleWithBackButton from "../../../../components/TitleWithBackButton";

import useAxios from "../../../../hooks/useAxios";
import { CategoryListItem, Response } from "../../../../types/response";
import {
  ADMIN_EDIT_CATEGORY_TITLE,
  API_ADMIN_CATEGORY,
  API_METHOD_GET,
  API_METHOD_PUT,
  PATH_ADMIN_CATEGORY,
} from "../../../../const/const";
import { CapitalizeFirstLetter } from "../../../../utils/formatter";
import ServerError from "../../../ServerError";
import PageLoader from "../../../../components/PageLoader";
import useTitle from "../../../../hooks/useTitle";

export default function EditCategory() {
  const navigate = useNavigate();
  const [category, setCategory] = useState("");
  const [showSaveModal, setShowSaveModal] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  useTitle(ADMIN_EDIT_CATEGORY_TITLE)

  const params = useParams();

  const {
    data,
    error: errorData,
    isLoading: isDataLoading,
  } = useAxios<Response<CategoryListItem>>(
    API_ADMIN_CATEGORY + "/" + params.id,
    API_METHOD_GET
  );

  const { error, isLoading, fetchData } = useAxios<Response<undefined>>(
    API_ADMIN_CATEGORY,
    API_METHOD_PUT
  );

  useEffect(() => {
    if (data?.data) {
      setCategory(data.data.name);
    }
  }, [data]);

  function handleSubmit() {
    fetchData({
      id: parseInt(params.id ? params.id : ""),
      name: category,
    });

    setShowResult(true);
    setShowSaveModal(false);
  }

  if (errorData && errorData[0].field === "server") {
    return <ServerError />;
  }

  if (isDataLoading) {
    return <PageLoader />;
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-16 w-[90%] max-w-[1108px] flex flex-col gap-8">
        <TitleWithBackButton>Edit Category</TitleWithBackButton>
        <form className="flex flex-col gap-8">
          <div className="w-full">
            <div className="w-full flex flex-col gap-8">
              <div className="w-full flex flex-col gap-1.5">
                <h1 className="text-primary-black font-semibold">Name</h1>
                <Input
                  value={category}
                  type="text"
                  name="category"
                  placeholder="Product Category"
                  onChange={(e) => {
                    setCategory(e.target.value);
                  }}
                />
              </div>
            </div>
          </div>
          <div className="w-full flex justify-end">
            <div className="w-[367px] flex justify-between gap-4">
              <Button
                submit={false}
                onClick={() => {
                  setShowSaveModal(true);
                }}
                disabled={isLoading || !category}
              >
                Save
              </Button>
              <Button
                submit={false}
                type="ghost-green"
                onClick={() => {
                  navigate(PATH_ADMIN_CATEGORY);
                }}
              >
                Cancel
              </Button>
            </div>
          </div>
        </form>
      </div>
      {showSaveModal &&
        createPortal(
          <ConfirmBox
            type="update"
            onYes={handleSubmit}
            onCancel={() => setShowSaveModal(false)}
          />,
          document.body
        )}
      {showResult &&
        createPortal(
          isLoading ? (
            <LoaderBox />
          ) : error ? (
            <ErrorBox onClose={() => setShowResult(false)}>
              <p className="font-bold text-primary-black text-4xl">Oops...</p>
              <p className="text-primary-black text-xl">
                {error && error[0].field === "server"
                  ? "Something is wrong, please try again"
                  : CapitalizeFirstLetter(error[0].detail)}
              </p>
            </ErrorBox>
          ) : (
            <SuccessBox onClose={() => navigate(PATH_ADMIN_CATEGORY)}>
              Category changed successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
