import { RouterProvider } from "react-router-dom";
import { Provider } from "react-redux";

import useRouter from "./router/router";
import { setupStore } from "./store/store";

const store = setupStore();

function App() {
  const router = useRouter();

  return (
    <Provider store={store}>
      <RouterProvider router={router} />
    </Provider>
  );
}

export default App;
