import { useNavigate } from "react-router-dom";
import { Plus } from "lucide-react";

import Button from "../../components/Button";
import AddressCard from "../../components/AddressCard";
import ErrorCard from "../../components/ErrorCard";
import PageLoader from "../../components/PageLoader";
import TitleWithBackButton from "../../components/TitleWithBackButton";
import ServerError from "../ServerError";

import useAxios from "../../hooks/useAxios";
import { UserAddress, Response } from "../../types/response";
import {
  API_ADDRESS_USER,
  API_METHOD_GET,
  ADDRESS_TITLE,
  PATH_CREATE_ADDRESS,
  PATH_EDIT_ADDRESS_EMPTY,
} from "../../const/const";
import useTitle from "../../hooks/useTitle";

export default function Addresses() {
  const navigate = useNavigate();

  useTitle(ADDRESS_TITLE);

  const { data, isLoading, error } = useAxios<Response<UserAddress[]>>(
    API_ADDRESS_USER,
    API_METHOD_GET
  );

  if (isLoading) {
    return <PageLoader />;
  }

  if (error && error[0].field === "server") {
    return <ServerError />;
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-7 md:my-16 w-[90%] max-w-[627px] flex flex-col gap-7 md:gap-10">
        <TitleWithBackButton>Address</TitleWithBackButton>
        <div className="w-full flex flex-col gap-5">
          {error && <ErrorCard errors={error} />}
          {data?.data &&
            data.data.map((addr) => {
              return (
                <AddressCard
                  key={addr.id}
                  address={addr}
                  onSelect={() => {
                    navigate(PATH_EDIT_ADDRESS_EMPTY + addr.id);
                  }}
                />
              );
            })}
        </div>
        <div className="flex justify-center">
          <div className="w-full max-w-[367px] flex justify-center gap-4">
            <Button
              onClick={() => {
                navigate(PATH_CREATE_ADDRESS);
              }}
            >
              <div className="flex justify-center gap-1.5">
                <Plus />
                <h1>Add a new address</h1>
              </div>
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
