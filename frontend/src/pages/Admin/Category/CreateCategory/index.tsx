import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createPortal } from "react-dom";

import Input from "../../../../components/Input";
import Button from "../../../../components/Button";
import ConfirmBox from "../../../../components/ConfirmBox";
import LoaderBox from "../../../../components/LoaderBox";
import ErrorBox from "../../../../components/ErrorBox";
import SuccessBox from "../../../../components/SuccessBox";
import TitleWithBackButton from "../../../../components/TitleWithBackButton";

import useAxios from "../../../../hooks/useAxios";
import { CapitalizeFirstLetter } from "../../../../utils/formatter";
import { Response } from "../../../../types/response";
import {
  ADMIN_ADD_CATEGORY_TITLE,
  API_ADMIN_CATEGORY,
  API_METHOD_POST,
  PATH_ADMIN_CATEGORY,
} from "../../../../const/const";
import useTitle from "../../../../hooks/useTitle";

export default function CreateCategory() {
  const navigate = useNavigate();
  const [category, setCategory] = useState("");
  const [isCategoryValid, setIsCategoryValid] = useState<boolean>(true);
  const [showSaveModal, setShowSaveModal] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  const { error, isLoading, fetchData } = useAxios<Response<undefined>>(
    API_ADMIN_CATEGORY,
    API_METHOD_POST
  );

  useTitle(ADMIN_ADD_CATEGORY_TITLE);

  function handleSubmit() {
    fetchData({ name: category });

    setShowResult(true);
    setShowSaveModal(false);
  }

  function handleCategoryChange(e: React.ChangeEvent<HTMLInputElement>) {
    setCategory(e.target.value);
    setIsCategoryValid(e.target.value.length > 2);
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <form className="my-16 w-[90%] max-w-[1108px] flex flex-col items-center gap-8">
        <div className="w-full flex justify-between gap-4">
          <TitleWithBackButton>Add Category</TitleWithBackButton>
        </div>
        <div className="w-full">
          <div className="w-full flex flex-col gap-8">
            <div className="relative w-full flex flex-col gap-1.5">
              <h1 className="text-primary-black font-semibold">Name</h1>
              <Input
                type="text"
                name="category"
                placeholder="Category name"
                valid={isCategoryValid}
                onChange={handleCategoryChange}
              />
              {!isCategoryValid && (
                <p className="absolute bottom-[-26px] text-sm text-primary-red pl-1 pt-1">
                  Name should be at least 3 characters
                </p>
              )}
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
      {showSaveModal &&
        createPortal(
          <ConfirmBox
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
              Category added successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
