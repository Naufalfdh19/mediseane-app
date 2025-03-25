import { useEffect, useState } from "react";

import Pagination from "../../../components/Pagination";
import TableHeader, {
  searchByType,
  sortByType,
} from "../../../components/TableHeader";

import {
  ADMIN_USERS_TITLE,
  API_ADMIN_USER,
  API_METHOD_GET,
} from "../../../const/const";
import { Response, UserData } from "../../../types/response";
import useAxios from "../../../hooks/useAxios";
import useTitle from "../../../hooks/useTitle";

const searchByList: searchByType[] = [
  {
    id: "1",
    name: "name",
  },
  {
    id: "2",
    name: "email",
  },
  {
    id: "3",
    name: "role",
  },
];

const sortByList: sortByType[] = [
  {
    id: "1",
    name: "created_at",
  },
  {
    id: "2",
    name: "name",
  },
];

export default function GetUser() {
  const limitPerPage = 10;
  const [search, setSearch] = useState("");
  const [searchBy, setSearchBy] = useState<searchByType>(searchByList[0]);
  const [filter, setFilter] = useState("");
  const [sortBy, setSortBy] = useState<sortByType>();
  const [asc, setAsc] = useState(false);

  const [page, setPage] = useState(1);
  const [totalPage, setTotalPage] = useState(1);

  const { data } = useAxios<Response<UserData>>(
    API_ADMIN_USER +
      `?page=${page}` +
      `&order=${asc ? "asc" : "desc"}` +
      (sortBy ? `&sort_by=${sortBy.name}` : "") +
      (searchBy && filter ? `&${searchBy.name}=${filter}` : "") +
      `&limit=${limitPerPage}`,
    API_METHOD_GET
  );

  useTitle(ADMIN_USERS_TITLE);

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

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-16 w-[90%] max-w-[1108px] flex flex-col items-center gap-8">
        <div className="w-full flex justify-between gap-4">
          <h1 className="text-center text-2xl font-bold text-primary-black">
            User
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
          readOnly={true}
        />

        <div className="w-full flex flex-col gap-8 bg-clip-content overflow-x-auto">
          <table className="min-w-[627px]">
            <thead>
              <tr className="border-b border-b-primary-gray h-[58px] text-primary-black">
                <th className="w-[60px]">No</th>
                <th className="text-left">Name</th>
                <th className="text-left">Email</th>
                <th className="text-left">Role</th>
              </tr>
            </thead>
            <tbody>
              {data?.data &&
              data.data.user_list &&
              data.data.user_list.length > 0
                ? data.data.user_list.map((item, idx) => {
                    return (
                      <tr className="h-[58px] text-primary-black" key={item.id}>
                        <td className="text-center w-[60px]">
                          {limitPerPage * (page - 1) + idx + 1}
                        </td>
                        <td className="text-left">{item.name}</td>
                        <td className="text-left">{item.email}</td>
                        <td
                          className={`text-left ${
                            item.role.toLocaleLowerCase() === "admin"
                              ? "text-primary-green"
                              : ""
                          }`}
                        >
                          {item.role}
                        </td>
                      </tr>
                    );
                  })
                : undefined}
            </tbody>
          </table>
          {!(
            data?.data &&
            data.data.user_list &&
            data.data.user_list.length > 0
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
    </div>
  );
}
