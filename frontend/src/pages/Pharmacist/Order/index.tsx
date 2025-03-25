import { useEffect, useState } from "react";
import { Ban, Send } from "lucide-react";
import { createPortal } from "react-dom";

import Pagination from "../../../components/Pagination";
import TableHeader, {
  searchByType,
  sortByType,
} from "../../../components/TableHeader";
import Button from "../../../components/Button";
import PageLoader from "../../../components/PageLoader";
import LoaderBox from "../../../components/LoaderBox";
import ErrorBox from "../../../components/ErrorBox";
import SuccessBox from "../../../components/SuccessBox";
import ConfirmBox from "../../../components/ConfirmBox";

import useAxios from "../../../hooks/useAxios";
import { OrdersData, OrdersListItem, Response } from "../../../types/response";
import {
  API_METHOD_DELETE,
  API_METHOD_GET,
  API_METHOD_PATCH,
  API_PHARMACIST_ORDERS,
  ORDER_STATUS_CANCELED,
  ORDER_STATUS_PENDING,
  ORDER_STATUS_PROCESSING,
  ORDER_STATUS_SHIPPED,
  PHARMACIST_ORDERS_TITLE,
} from "../../../const/const";
import { parseTimestamp } from "../../../utils/timestamp";
import Modal from "../../../components/Modal/Modal";
import { CapitalizeFirstLetter } from "../../../utils/formatter";
import { statusColor } from "../../../utils/orders";
import useTitle from "../../../hooks/useTitle";

const searchByList: searchByType[] = [
  {
    id: "1",
    name: "status",
  },
];

const sortByList: sortByType[] = [
  {
    id: "1",
    name: "status",
  },
  {
    id: "2",
    name: "created_at",
  },
];

export default function PharmacistOrders() {
  const limitPerPage = 10;
  const [search, setSearch] = useState("");
  const [page, setPage] = useState<number>(1);
  const [searchBy, setSearchBy] = useState<searchByType>(searchByList[0]);
  const [filter, setFilter] = useState("");
  const [sortBy, setSortBy] = useState<sortByType>();
  const [totalPage, setTotalPage] = useState<number>(1);
  const [asc, setAsc] = useState(true);
  const [productsId, setProductsId] = useState<number>();
  const [updateId, setUpdateId] = useState<number>();
  const [cancelId, setCancelId] = useState<number>();
  const [showProducts, setShowProducts] = useState<boolean>(false);
  const [showUpdateConfirm, setShowUpdateConfirm] = useState<boolean>(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState<boolean>(false);
  const [showResult, setShowResult] = useState<boolean>(false);

  const { data, isLoading, fetchData } = useAxios<Response<OrdersData>>(
    API_PHARMACIST_ORDERS +
      `?page=${page}` +
      `&limit=${limitPerPage}` +
      `&order=${asc ? "asc" : "desc"}` +
      (sortBy ? `&sort_by=${sortBy.name}` : "") +
      (searchBy && filter ? `&${searchBy.name}=${filter}` : "")
  );

  const {
    data: products,
    isLoading: isLoadingProducts,
    error: errProducts,
    fetchData: fetchProducts,
  } = useAxios<Response<OrdersListItem>>(
    API_PHARMACIST_ORDERS + "/" + productsId,
    API_METHOD_GET,
    false,
    true
  );

  const {
    isLoading: isLoadingUpdate,
    error: errUpdate,
    fetchData: fetchUpdate,
  } = useAxios<Response<undefined>>(
    API_PHARMACIST_ORDERS + "/" + updateId,
    API_METHOD_PATCH
  );

  const {
    isLoading: isLoadingCancel,
    error: errCancel,
    fetchData: fetchCancel,
  } = useAxios<Response<undefined>>(
    API_PHARMACIST_ORDERS + "/" + cancelId,
    API_METHOD_DELETE
  );

  useTitle(PHARMACIST_ORDERS_TITLE);

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

  useEffect(() => {
    if (showProducts && productsId) {
      fetchProducts();
    }
  }, [showProducts]);

  function handleProductDetails(order_id: number) {
    setProductsId(order_id);
    setShowProducts(true);
  }

  function handleUpdate() {
    fetchUpdate({
      status: "Shipped",
    });
    setShowResult(true);
    setShowUpdateConfirm(false);
  }

  function handleCancel() {
    fetchCancel();
    setShowResult(true);
    setShowDeleteConfirm(false);
  }

  function handleClose() {
    setShowResult(false);
    setUpdateId(undefined);
    setCancelId(undefined);
    fetchData();
  }

  if (isLoading) {
    return <PageLoader />;
  }

  return (
    <div className="grow flex bg-primary-white justify-center">
      <div className="my-16 w-[90%] max-w-[1108px] flex flex-col items-center gap-8">
        <div className="w-full flex justify-between gap-4">
          <h1 className="text-center text-2xl font-bold text-primary-black">
            Orders
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
                <th className="w-[250px] text-left">Order ID</th>
                <th className="w-[200px] text-left">Status</th>
                <th className="w-[250px] text-left">Order Date</th>
                <th className="text-center">Action</th>
              </tr>
            </thead>
            <tbody>
              {data?.data && data.data.orders && data.data.orders.length > 0
                ? data.data.orders.map((item, idx) => {
                    return (
                      <tr
                        key={item.order_id}
                        className="h-[58px] text-primary-black"
                      >
                        <td className="text-center w-[60px]">
                          {limitPerPage * (page - 1) + idx + 1}
                        </td>
                        <td
                          onClick={() => {
                            handleProductDetails(item.order_id);
                          }}
                          className="text-left text-blue font-bold underline cursor-pointer"
                        >
                          ORDER-{item.order_id}
                        </td>
                        <td className={`text-left ${statusColor(item.status)}`}>
                          {item.status}
                        </td>
                        <td className="text-left">
                          {parseTimestamp(item.created_at)}
                        </td>
                        <td className="text-center w-[151px]">
                          <div className="flex justify-center gap-2">
                            {item.status.toLowerCase() ===
                              ORDER_STATUS_PROCESSING && (
                              <div className="w-10">
                                <Button
                                  onClick={() => {
                                    setUpdateId(item.order_id);
                                    setShowUpdateConfirm(true);
                                  }}
                                  square={true}
                                  size="xs"
                                >
                                  <div
                                    aria-label="send-btn"
                                    className="w-full flex justify-center items-center"
                                  >
                                    <Send className="w-[13px] h-[13px]" />
                                  </div>
                                </Button>
                              </div>
                            )}
                            {(item.status.toLowerCase() ===
                              ORDER_STATUS_PENDING ||
                              item.status.toLowerCase() ===
                                ORDER_STATUS_PROCESSING) && (
                              <div className="w-10">
                                <Button
                                  onClick={() => {
                                    setCancelId(item.order_id);
                                    setShowDeleteConfirm(true);
                                  }}
                                  square={true}
                                  size="xs"
                                  type="ghost-red"
                                >
                                  <div
                                    aria-label="cancel-btn"
                                    className="w-full flex justify-center items-center"
                                  >
                                    <Ban className="w-[11px] h-[14px]" />
                                  </div>
                                </Button>
                              </div>
                            )}
                            {!(
                              item.status.toLowerCase() ===
                                ORDER_STATUS_PENDING ||
                              item.status.toLowerCase() ===
                                ORDER_STATUS_PROCESSING
                            ) && <p className="text-primary-black">-</p>}
                          </div>
                        </td>
                      </tr>
                    );
                  })
                : undefined}
            </tbody>
          </table>
          {!(data?.data && data.data.orders && data.data.orders.length > 0) && (
            <div className="w-full flex justify-center">
              <p className="text-primary-gray">Item Not Found</p>
            </div>
          )}
        </div>
        <div className="w-full flex justify-center md:justify-end">
          <Pagination page={page} setPage={setPage} totalPage={totalPage} />
        </div>
      </div>
      {showProducts &&
        createPortal(
          isLoadingProducts ? (
            <LoaderBox />
          ) : errProducts ? (
            <ErrorBox onClose={() => setShowProducts(false)}>
              <p className="font-bold text-primary-black text-4xl">Oops...</p>
              <p className="text-primary-black text-xl">
                Something is wrong, please try again
              </p>
            </ErrorBox>
          ) : (
            <Modal
              onClose={() => {
                setShowProducts(false);
              }}
            >
              <div className="w-full py-8 px-2 flex flex-col items-center justify-center gap-6">
                <h1 className="w-[95%] font-semibold text-primary-black text-2xl text-center">
                  Order ID: {productsId}
                </h1>
                <div className="overflow-auto w-full max-h-[60vh]">
                  <table className="min-w-[600px]">
                    <thead>
                      <tr className="border-b border-b-primary-gray h-[58px] text-primary-black">
                        <th className="text-center">ID</th>
                        <th className="text-left">Name</th>
                        <th className="text-left">Quantity</th>
                        <th className="text-center">Image</th>
                      </tr>
                    </thead>
                    <tbody>
                      {products?.data && products.data.product_list
                        ? products.data.product_list.map((item) => {
                            return (
                              <tr
                                key={item.product_id}
                                className="h-[58px] text-primary-black"
                              >
                                <td className="text-center w-[60px]">
                                  {item.product_id}
                                </td>
                                <td className="text-left">{item.name}</td>
                                <td className="text-left">{item.quantity}</td>
                                <td className="flex justify-center items-center">
                                  <img
                                    className="w-[40px] h-[40px]"
                                    src={item.image}
                                    alt={item.name}
                                  />
                                </td>
                              </tr>
                            );
                          })
                        : undefined}
                    </tbody>
                  </table>
                  {!(products?.data && products.data.product_list) && (
                    <div className="w-full flex justify-center">
                      <p className="text-primary-gray">Item Not Found</p>
                    </div>
                  )}
                </div>
                <div className="w-[130px] md:w-[163px]">
                  <Button
                    submit={false}
                    type="ghost"
                    size="md"
                    onClick={() => {
                      setShowProducts(false);
                    }}
                  >
                    Close
                  </Button>
                </div>
              </div>
            </Modal>
          ),
          document.body
        )}
      {showUpdateConfirm &&
        createPortal(
          <ConfirmBox
            type="ship"
            onYes={handleUpdate}
            onCancel={() => {
              setUpdateId(undefined);
              setShowUpdateConfirm(false);
            }}
          />,
          document.body
        )}
      {showDeleteConfirm &&
        createPortal(
          <ConfirmBox
            type="cancel"
            onYes={handleCancel}
            onCancel={() => {
              setCancelId(undefined);
              setShowDeleteConfirm(false);
            }}
          />,
          document.body
        )}
      {showResult &&
        createPortal(
          isLoadingUpdate || isLoadingCancel ? (
            <LoaderBox />
          ) : errUpdate || errCancel ? (
            <ErrorBox onClose={() => setShowResult(false)}>
              <p className="font-bold text-primary-black text-4xl">Oops...</p>
              <p className="text-primary-black text-xl">
                {(errUpdate && errUpdate[0].field === "server") ||
                (errCancel && errCancel[0].field === "server")
                  ? "Something is wrong, please try again"
                  : (errUpdate && CapitalizeFirstLetter(errUpdate[0].detail)) ||
                    (errCancel && CapitalizeFirstLetter(errCancel[0].detail))}
              </p>
            </ErrorBox>
          ) : (
            <SuccessBox onClose={handleClose}>
              Order {cancelId ? ORDER_STATUS_CANCELED : ORDER_STATUS_SHIPPED}{" "}
              successfully!
            </SuccessBox>
          ),
          document.body
        )}
    </div>
  );
}
