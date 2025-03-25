import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { createPortal } from "react-dom";
import { Pencil, Trash } from "lucide-react";

import Pagination from "../../../components/Pagination";
import TableHeader, {
  searchByType,
  sortByType,
} from "../../../components/TableHeader";
import Button from "../../../components/Button";
import ConfirmBox from "../../../components/ConfirmBox";
import PageLoader from "../../../components/PageLoader";
import LoaderBox from "../../../components/LoaderBox";
import ErrorBox from "../../../components/ErrorBox";
import SuccessBox from "../../../components/SuccessBox";

import useAxios from "../../../hooks/useAxios";
import { CapitalizeFirstLetter } from "../../../utils/formatter";
import { PartnersData, Response } from "../../../types/response";
import {
  ADMIN_PARTNERS_TITLE,
  API_ADMIN_PARTNERS,
  API_METHOD_DELETE,
  PATH_ADMIN_ADD_PARTNERS,
  PATH_ADMIN_EDIT_PARTNERS_EMPTY,
} from "../../../const/const";
import useTitle from "../../../hooks/useTitle";

const searchByList: searchByType[] = [
  {
    id: "1",
    name: "name",
  },
  {
    id: "2",
    name: "year_founded",
  },
  {
    id: "3",
    name: "active_days",
  },
];

const sortByList: sortByType[] = [
  {
    id: "1",
    name: "created_at",
  },
  {
    id: "2",
    name: "year_founded",
  },
];

export default function AdminPartners() {
  const limitPerPage = 10;

  const [search, setSearch] = useState("");
  const [page, setPage] = useState<number>(1);
  const [searchBy, setSearchBy] = useState<searchByType>(searchByList[0]);
  const [filter, setFilter] = useState("");
  const [sortBy, setSortBy] = useState<sortByType>();
  const [totalPage, setTotalPage] = useState<number>(1);
  const [asc, setAsc] = useState(false);

  const [deletedID, setDeletedID] = useState<number>();
  const [showDeleteConfirm, setShowDeleteConfirm] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  const { data, isLoading } = useAxios<Response<PartnersData>>(
    API_ADMIN_PARTNERS +
      `?page=${page}` +
      `&order=${asc ? "asc" : "desc"}` +
      (sortBy ? `&sort_by=${sortBy.name}` : "") +
      (searchBy && filter ? `&${searchBy.name}=${filter}` : "")
  );

  const {
    error,
    isLoading: isDeleteLoading,
    fetchData: fetchDelete,
  } = useAxios<Response<PartnersData>>(
    API_ADMIN_PARTNERS + "/" + deletedID,
    API_METHOD_DELETE
  );

  useTitle(ADMIN_PARTNERS_TITLE);

  const navigate = useNavigate();

  useEffect(() => {
    const currPage =
      data &&
      data.data &&
      data.data.pagination &&
      data.data.pagination.current_page >= 1
        ? data.data.pagination.current_page
        : 1;
    const totalPage =
      data &&
      data.data &&
      data.data.pagination &&
      data.data.pagination.total_page >= 1
        ? data.data.pagination.total_page
        : 1;

    setTotalPage(totalPage);
    setPage(currPage);
  }, [data]);

  useEffect(() => {
    if (page < 1) {
      setPage(1);
    }
    if (page > totalPage) {
      setPage(totalPage);
    }
  }, [page, totalPage]);

  function handleDelete() {
    fetchDelete();
    setShowResult(true);
  }

  if (isLoading) {
    return <PageLoader />;
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-16 w-[90%] max-w-[1108px] flex flex-col items-center gap-8">
        <div className="w-full flex justify-between gap-4">
          <h1 className="text-center text-2xl font-bold text-primary-black">
            Partners
          </h1>
        </div>
        <TableHeader
          search={search}
          searchBy={searchBy}
          setSearchBy={setSearchBy}
          searchByList={searchByList}
          setSearch={setSearch}
          setFilter={setFilter}
          sortBy={sortBy}
          setSortBy={setSortBy}
          sortByList={sortByList}
          asc={asc}
          setAsc={setAsc}
          onAdd={() => navigate(PATH_ADMIN_ADD_PARTNERS)}
        />
        <div className="w-full flex flex-col gap-8 bg-clip-content overflow-x-auto">
          <table className="min-w-[627px]">
            <thead>
              <tr className="border-b border-b-primary-gray h-[58px] text-primary-black">
                <th className="w-[60px] px-4">No</th>
                <th className="w-[250px] text-left px-4">Name</th>
                <th className="w-[200px] text-left px-4">Year founded</th>
                <th className="w-[250px] text-left px-4">Active days</th>
                <th className="w-[200px] text-left px-4">Operational hour</th>
                <th className="text-left px-4">Active</th>
                <th className="text-center px-4">Action</th>
              </tr>
            </thead>
            <tbody>
              {data?.data &&
              data.data.partner_list &&
              data.data.partner_list.length > 0
                ? data.data.partner_list.map((item, idx) => {
                    return (
                      <tr key={item.id} className="h-[58px] text-primary-black">
                        <td className="text-center w-[60px] px-4">
                          {limitPerPage * (page - 1) + idx + 1}
                        </td>
                        <td className="text-left px-4">{item.name}</td>
                        <td className="text-left px-4">{item.year_founded}</td>
                        <td className="text-left px-4">{item.active_days}</td>
                        <td className="text-left px-4">
                          {item.start_hour + "-" + item.end_hour}
                        </td>
                        <td className="text-left px-4">
                          {item.is_active ? "Yes" : "No"}
                        </td>
                        <td className="text-center w-[151px] px-4">
                          <div className="flex justify-center gap-2">
                            <div className="w-10">
                              <Button
                                onClick={() => {
                                  navigate(
                                    PATH_ADMIN_EDIT_PARTNERS_EMPTY + item.id
                                  );
                                }}
                                square={true}
                                size="xs"
                              >
                                <div className="w-full flex justify-center items-center">
                                  <Pencil
                                    className="w-[13px] h-[13px]"
                                    aria-label="Edit"
                                  />
                                </div>
                              </Button>
                            </div>
                            <div className="w-10">
                              <Button
                                onClick={() => {
                                  setDeletedID(item.id);
                                  setShowDeleteConfirm(true);
                                }}
                                square={true}
                                size="xs"
                                type="ghost-red"
                              >
                                <div className="w-full flex justify-center items-center">
                                  <Trash
                                    className="w-[11px] h-[14px]"
                                    aria-label="Delete"
                                  />
                                </div>
                              </Button>
                            </div>
                          </div>
                        </td>
                      </tr>
                    );
                  })
                : undefined}
            </tbody>
          </table>
          {!(
            data?.data &&
            data.data.partner_list &&
            data.data.partner_list.length > 0
          ) && (
            <div className="w-full flex justify-center">
              <p className="text-primary-gray">Item Not Found</p>
            </div>
          )}
        </div>
        <div className="w-full flex justify-center md:justify-end">
          <Pagination page={page} setPage={setPage} totalPage={totalPage} />
        </div>
      </div>
      {showDeleteConfirm &&
        createPortal(
          <ConfirmBox
            type="delete"
            onYes={() => {
              handleDelete();
              setShowDeleteConfirm(false);
            }}
            onCancel={() => {
              setShowDeleteConfirm(false);
            }}
          />,
          document.body
        )}
      {showResult &&
        createPortal(
          isDeleteLoading ? (
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
            <SuccessBox onClose={() => window.location.reload()}>
              Partner deleted successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
