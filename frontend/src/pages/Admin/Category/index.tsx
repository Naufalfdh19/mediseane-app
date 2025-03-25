import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { createPortal } from "react-dom";
import { Pencil, Trash } from "lucide-react";

import TableHeader, { sortByType } from "../../../components/TableHeader";
import Button from "../../../components/Button";
import Pagination from "../../../components/Pagination";
import ConfirmBox from "../../../components/ConfirmBox";
import LoaderBox from "../../../components/LoaderBox";
import ErrorBox from "../../../components/ErrorBox";
import SuccessBox from "../../../components/SuccessBox";

import useAxios from "../../../hooks/useAxios";
import { CategoryData, Response } from "../../../types/response";
import {
  ADMIN_CATEGORY_TITLE,
  API_ADMIN_CATEGORY,
  API_CATEGORIES,
  API_METHOD_DELETE,
  API_METHOD_GET,
  PATH_ADMIN_CREATE_CATEGORY,
  PATH_ADMIN_EDIT_CATEGORY_EMPTY,
} from "../../../const/const";
import { CapitalizeFirstLetter } from "../../../utils/formatter";
import useTitle from "../../../hooks/useTitle";

const sortByList: sortByType[] = [
  {
    id: "1",
    name: "Date",
  },
  {
    id: "2",
    name: "Name",
  },
];

export default function Category() {
  const limitPerPage = 10;
  const [search, setSearch] = useState("");
  const [filter, setFilter] = useState("");
  const [sortBy, setSortBy] = useState<sortByType>();
  const [asc, setAsc] = useState(false);
  const [deleteId, setDeleteId] = useState<number>();
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [totalPage, setTotalPage] = useState(1);

  useTitle(ADMIN_CATEGORY_TITLE);

  const { data, fetchData } = useAxios<Response<CategoryData>>(
    API_CATEGORIES +
      `?page=${page}` +
      `&limit=${limitPerPage}` +
      `&order=${asc ? "asc" : "desc"}` +
      (sortBy ? `&sortBy=${sortBy.name}` : "") +
      (filter ? `&filter=${filter}` : ""),
    API_METHOD_GET
  );

  const {
    error,
    isLoading,
    fetchData: fetchDelete,
  } = useAxios<Response<undefined>>(
    API_ADMIN_CATEGORY + "/" + deleteId,
    API_METHOD_DELETE
  );

  useEffect(() => {
    const currPage =
      data && data.data && data.data.pagination
        ? data.data.pagination.current_page
        : 1;
    const totalPage =
      data && data.data && data.data.pagination
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

  function parseTimestamp(create_timestamp: string) {
    const createdAt = new Date(create_timestamp);
    createdAt.setHours(createdAt.getHours() - 7);
    const createdDate = createdAt.toLocaleDateString("en-US");
    const createdTime = createdAt.toLocaleTimeString("en-US");

    return `${createdDate} at ${createdTime}`;
  }

  function handleDelete() {
    fetchDelete();

    setShowResult(true);
    setShowDeleteConfirm(false);
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-16 w-[90%] max-w-[1108px] flex flex-col items-center gap-8">
        <div className="w-full flex justify-between gap-4">
          <h1 className="text-center text-2xl font-bold text-primary-black">
            Category
          </h1>
        </div>
        <TableHeader
          search={search}
          setSearch={setSearch}
          setFilter={setFilter}
          sortBy={sortBy}
          setSortBy={setSortBy}
          sortByList={sortByList}
          asc={asc}
          setAsc={setAsc}
          onAdd={() => {
            navigate(PATH_ADMIN_CREATE_CATEGORY);
          }}
        />

        <div className="w-full flex flex-col gap-8 bg-clip-content overflow-x-auto">
          <table className="min-w-[627px]">
            <thead>
              <tr className="border-b border-b-primary-gray h-[58px]">
                <th className="w-[60px]">No</th>
                <th className="text-left">Name</th>
                <th className="text-left">Date</th>
                <th className="w-[151px]">Action</th>
              </tr>
            </thead>
            <tbody>
              {data?.data && data.data.list_item.length > 0
                ? data.data.list_item.map((item, idx) => {
                    return (
                      <tr className="h-[58px]" key={item.id}>
                        <td className="text-center w-[60px]">
                          {limitPerPage * (page - 1) + idx + 1}
                        </td>
                        <td className="text-left">{item.name}</td>
                        <td className="text-left">
                          {parseTimestamp(item.updated_at)}
                        </td>
                        <td className="text-center w-[151px]">
                          <div className="flex justify-center gap-2">
                            <div className="w-10">
                              <Button
                                onClick={() => {
                                  navigate(
                                    PATH_ADMIN_EDIT_CATEGORY_EMPTY + item.id
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
                                  setDeleteId(item.id);
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

          {data?.data && data.data.list_item.length === 0 && (
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
            onYes={handleDelete}
            onCancel={() => setShowDeleteConfirm(false)}
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
            <SuccessBox
              onClose={() => {
                fetchData();
                setShowResult(false);
              }}
            >
              Category deleted successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
