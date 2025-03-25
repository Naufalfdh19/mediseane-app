import { useState } from "react";
import { Hospital, ShoppingCart } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { createPortal } from "react-dom";

import AddToCartButton from "../AddToCartButton";
import Button from "../Button";
import Modal from "../Modal/Modal";

import { ProductListItem } from "../../types/response";
import { useUserState } from "../../store/user";
import {
  PATH_LOGIN,
  PATH_PRODUCT_DETAILS_EMPTY,
  PATH_PROFILE,
} from "../../const/const";

interface ProductCardProps {
  ref?: (node: Element | null) => void;
  product: ProductListItem;
}

export default function ProductCard(props: ProductCardProps) {
  const [isNotAbleToShop, setIsNotAbleToShop] = useState<boolean>(false);

  const navigate = useNavigate();
  const userState = useUserState();

  function handleNotLogin() {
    navigate(PATH_LOGIN);
  }

  function handleNotVerified() {
    navigate(PATH_PROFILE);
  }

  return (
    <div
      ref={props.ref}
      onClick={() => {
        navigate(
          PATH_PRODUCT_DETAILS_EMPTY + props.product.pharmacy_product_id
        );
      }}
      className="bg-pure-white w-full overflow-hidden cursor-pointer shadow-xl border border-primary-gray/30 rounded-lg"
    >
      <div className="w-full flex flex-col items-center gap-2">
        <div className=" relative w-full h-[148px] overflow-hidden flex justify-center items-center">
          <img
            className="object-cover w-full"
            src={props.product.image.split("'")[1]}
            alt={props.product.name}
          />
        </div>
        <div className="w-full px-3.5 py-3 md:p-5 flex flex-col gap-5">
          <div className="flex flex-col gap-2.5">
            <div className="flex flex-col gap-1">
              <h1 className="text-primary-black text-lg truncate">
                {props.product.name}
              </h1>
              <h1 className="text-primary-gray text-base truncate">
                {props.product.manufacture}
              </h1>
              <div className="flex gap-1.5 items-center text-primary-gray ">
                <Hospital className="h-5" />
                <h1 className="text-base truncate">
                  {props.product.pharmacy_name}
                </h1>
              </div>
            </div>
            <h1 className="text-primary-green text-base font-semibold truncate">
              {props.product.price ? "Rp" + props.product.price : ""}
            </h1>
          </div>
          <div
            onClick={(e) => e.stopPropagation()}
            className="flex flex-col gap-2 justify-center"
          >
            {userState.id !== 0 && userState.is_verified ? (
              <AddToCartButton
                product_id={props.product.pharmacy_product_id}
                size="sm"
              />
            ) : (
              <Button onClick={() => setIsNotAbleToShop(true)} size="sm">
                <div className="flex justify-center items-center gap-1.5">
                  <ShoppingCart className="hidden md:inline" />
                  <p>Add to cart</p>
                </div>
              </Button>
            )}
            {isNotAbleToShop &&
              createPortal(
                <Modal onClose={() => setIsNotAbleToShop(false)}>
                  <div className="w-full py-8 px-2 flex flex-col items-center justify-center gap-6">
                    <ShoppingCart className="w-[100px] h-[100px] md:w-[130px] md:h-[130px] text-primary-green" />
                    <p className="px-2 font-semibold text-primary-black text-2xl md:text-3xl text-center">
                      {userState.id === 0 ? "Login" : "Verify email"} before
                      adding to cart
                    </p>
                    <div className="w-full flex flex-wrap gap-3 md:gap-7 justify-center items-center">
                      <div className="w-[120px] md:w-[163px]">
                        <Button
                          submit={false}
                          size="md"
                          onClick={
                            userState.id === 0
                              ? handleNotLogin
                              : handleNotVerified
                          }
                        >
                          {userState.id === 0 ? "Login now" : "Verify email"}
                        </Button>
                      </div>
                      <div className="w-[120px] md:w-[163px]">
                        <Button
                          submit={false}
                          type="ghost"
                          size="md"
                          onClick={() => setIsNotAbleToShop(false)}
                        >
                          Cancel
                        </Button>
                      </div>
                    </div>
                  </div>
                </Modal>,
                document.body
              )}
          </div>
        </div>
      </div>
    </div>
  );
}
