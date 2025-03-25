import { Link, useNavigate } from "react-router-dom";
import { useRef, useState } from "react";
import {
  Building2,
  ChartBarStacked,
  ChevronDown,
  CircleUserRound,
  FlaskConical,
  Hospital,
  LogOut,
  LucideComputer,
  Pill,
} from "lucide-react";

import {
  PATH_ADMIN_PARTNERS,
  PATH_ADMIN_PHARMACY,
  PATH_ADMIN_PHARMACIST,
  PATH_ADMIN_CATEGORY,
  PATH_ADMIN_DASHBOARD,
  PATH_ADMIN_LOGOUT,
  PATH_ADMIN_USERS,
  PATH_ADMIN_PRODUCTS,
} from "../../const/const";
import useWindowDimensions from "../../hooks/useWindowDimensions";
import useOutsideClick from "../../hooks/useOutsideClick";

const NavWidthBr = 1024;
const NavMenu = [
  { name: "Dashboard", symbol: <LucideComputer />, path: PATH_ADMIN_DASHBOARD },
  { name: "Users", symbol: <CircleUserRound />, path: PATH_ADMIN_USERS },
  {
    name: "Pharmacists",
    symbol: <FlaskConical />,
    path: PATH_ADMIN_PHARMACIST,
  },
  { name: "Pharmacies", symbol: <Hospital />, path: PATH_ADMIN_PHARMACY },
  { name: "Partners", symbol: <Building2 />, path: PATH_ADMIN_PARTNERS },
  {
    name: "Categories",
    symbol: <ChartBarStacked />,
    path: PATH_ADMIN_CATEGORY,
  },
  { name: "Products", symbol: <Pill />, path: PATH_ADMIN_PRODUCTS },
  { name: "Logout", symbol: <LogOut />, path: PATH_ADMIN_LOGOUT },
];

function AdminHeader() {
  const [isOpen, setIsOpen] = useState<boolean>(false);

  const navigate = useNavigate();
  const { width } = useWindowDimensions();

  const dropdownRef = useRef<HTMLDivElement>(null);
  useOutsideClick({
    ref: dropdownRef,
    handler: () => setIsOpen(false),
  });

  const dropdownClass =
    "absolute bg-primary-white w-full max-h-52 overflow-y-auto py-3 rounded shadow-md z-10 top-full left-0 mt-1";

  return (
    <div
      ref={dropdownRef}
      className="min-h-14 w-full bg-primary-white flex items-center border-b border-solid border-primary-gray/30"
    >
      {width > NavWidthBr ? (
        <div className="w-[90%] lg:w-[80%] mx-auto flex justify-between gap-4">
          <h1 className="text-base font-semibold text-primary-black">
            medi<span className="text-primary-green">SEA</span>ne
          </h1>
          <div className="flex gap-4">
            {NavMenu.map((menu) => {
              return (
                <Link
                  key={menu.name}
                  to={menu.path}
                  className={`flex gap-1 items-center text-sm font-semibold hover:text-primary-black ${
                    menu.path !== "" &&
                    window.location.pathname.indexOf(menu.path) === 0
                      ? "text-primary-green"
                      : "text-primary-gray"
                  }`}
                >
                  <div className="hidden xl:flex justify-center items-center">
                    {menu.symbol}
                  </div>
                  <p>{menu.name}</p>
                </Link>
              );
            })}
          </div>
        </div>
      ) : (
        <div className="w-full h-full relative flex justify-center items-center">
          <div className="w-[90%] lg:w-[80%] mx-auto flex justify-between">
            <h1
              className="text-base font-semibold text-primary-black cursor-pointer"
              onClick={() => {
                navigate(PATH_ADMIN_DASHBOARD);
              }}
            >
              medi<span className="text-primary-green">SEA</span>ne
            </h1>
            <button
              aria-label="Toggle dropdown"
              aria-haspopup="true"
              aria-expanded={isOpen}
              type="button"
              onClick={() => setIsOpen(!isOpen)}
              className={`flex justify-between items-center gap-5 rounded `}
            >
              <ChevronDown
                size={20}
                className={`transform duration-500 ease-in-out ${
                  isOpen ? "rotate-180" : ""
                }`}
              />
            </button>
          </div>
          {isOpen && (
            <div aria-label="Dropdown menu" className={dropdownClass}>
              <ul
                role="menu"
                aria-orientation="vertical"
                className="leading-10"
              >
                {NavMenu?.map((menu) => (
                  <li
                    onClick={() => {
                      navigate(menu.path);
                    }}
                    key={menu.name}
                    className={`flex items-center gap-1 cursor-pointer hover:bg-primary-gray hover:text-primary-white px-5  ${
                      window.location.pathname.indexOf(menu.path) === 0
                        ? "text-primary-green"
                        : "text-primary-black"
                    }`}
                  >
                    <div className="flex justify-center items-center">
                      {menu.symbol}
                    </div>
                    <span>{menu.name}</span>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default AdminHeader;
