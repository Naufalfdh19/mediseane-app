import { useEffect, useState } from "react";
import { Pencil, Trash } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { createPortal } from "react-dom";

import Pagination from "../../../components/Pagination";
import TableHeader, {
  searchByType,
  sortByType,
} from "../../../components/TableHeader";
import Button from "../../../components/Button";
import ConfirmBox from "../../../components/ConfirmBox";
import LoaderBox from "../../../components/LoaderBox";
import ErrorBox from "../../../components/ErrorBox";
import SuccessBox from "../../../components/SuccessBox";

import {
  ADMIN_PHARMACIES_TITLE,
  API_ADMIN_PHARMACY,
  API_METHOD_DELETE,
  API_METHOD_GET,
  PATH_ADMIN_CREATE_PHARMACY,
  PATH_ADMIN_EDIT_PHARMACY_EMPTY,
} from "../../../const/const";
import { PharmacyData, Response } from "../../../types/response";
import useAxios from "../../../hooks/useAxios";
import { CapitalizeFirstLetter } from "../../../utils/formatter";
import useTitle from "../../../hooks/useTitle";

const searchByList: searchByType[] = [
  {
    id: "1",
    name: "name",
  },
  {
    id: "2",
    name: "city",
  },
];

const sortByList: sortByType[] = [
  {
    id: "1",
    name: "last_update",
  },
];

export default function Pharmacy() {
  const limitPerPage = 10;
  const [search, setSearch] = useState("");
  const [searchBy, setSearchBy] = useState<searchByType>(searchByList[0]);
  const [filter, setFilter] = useState("");
  const [sortBy, setSortBy] = useState<sortByType>();
  const [asc, setAsc] = useState(false);
  const [deleteId, setDeleteId] = useState<number>();
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [totalPage, setTotalPage] = useState(1);

  const { data, fetchData } = useAxios<Response<PharmacyData>>(
    API_ADMIN_PHARMACY +
      `?page=${page}` +
      `&order=${asc ? "asc" : "desc"}` +
      (sortBy ? `&sort_by=${sortBy.name}` : "") +
      (searchBy && filter ? `&${searchBy.name}=${filter}` : "") +
      `&limit=${limitPerPage}`,
    API_METHOD_GET
  );

  const {
    error,
    isLoading,
    fetchData: fetchDelete,
  } = useAxios<Response<undefined>>(
    API_ADMIN_PHARMACY + "/" + deleteId,
    API_METHOD_DELETE
  );

  useTitle(ADMIN_PHARMACIES_TITLE);

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
    if (totalPage < 1) {
      setTotalPage(1);
    }
    if (page > totalPage) {
      setPage(totalPage);
    }
  }, [page, totalPage]);

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
            Pharmacy
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
          onAdd={() => {
            navigate(PATH_ADMIN_CREATE_PHARMACY);
          }}
        />

        <div className="w-full flex flex-col gap-8 bg-clip-content overflow-x-auto">
          <table className="min-w-[627px]">
            <thead>
              <tr className="border-b border-b-primary-gray h-[58px] text-primary-black">
                <th className="w-[60px] px-4">No</th>
                <th className="text-left px-4">Name</th>
                <th className="text-left px-4">Logo</th>
                <th className="text-left px-4">Partner</th>
                <th className="text-left px-4">Address</th>
                <th className="w-[151px] px-4">Action</th>
              </tr>
            </thead>
            <tbody>
              {data?.data &&
              data.data.pharmacies &&
              data.data.pharmacies.length > 0
                ? data.data.pharmacies.map((item, idx) => {
                    return (
                      <tr className="h-[58px] text-primary-black" key={item.id}>
                        <td className="text-center w-[60px] px-4">
                          {limitPerPage * (page - 1) + idx + 1}
                        </td>
                        <td className="text-left px-4">{item.name}</td>
                        <td className="text-left px-4">
                          {item.logo && (
                            <img
                              alt={item.name}
                              src={item.logo}
                              className="w-8"
                            />
                          )}
                        </td>
                        <td className="text-left px-4">{item.partner_name}</td>
                        <td className="text-left px-4">
                          {[
                            item.address,
                            item.sub_district.toUpperCase(),
                            item.district.toUpperCase(),
                            item.city.toUpperCase(),
                            item.province.toUpperCase(),
                            item.postal_code,
                          ].join(", ")}
                        </td>
                        <td className="text-center w-[151px]">
                          <div className="flex justify-center gap-2">
                            <div className="w-10">
                              <Button
                                onClick={() => {
                                  navigate(
                                    PATH_ADMIN_EDIT_PHARMACY_EMPTY + item.id
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
          {!(
            data?.data &&
            data.data.pharmacies &&
            data.data.pharmacies.length > 0
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
              Pharmacy deleted successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
