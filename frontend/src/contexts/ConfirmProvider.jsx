import { useCallback, useRef, useState } from "react";
import ConfirmDialog from "../components/common/ConfirmDialog";
import { ConfirmContext } from "./confirmContext";

const INITIAL_DIALOG = {
  open: false,
  title: "Xác nhận thao tác",
  message: "",
  confirmLabel: "Xác nhận",
  danger: false,
};

export default function ConfirmProvider({ children }) {
  const [dialog, setDialog] = useState(INITIAL_DIALOG);
  const resolverRef = useRef(null);

  const close = useCallback((result) => {
    resolverRef.current?.(result);
    resolverRef.current = null;
    setDialog(INITIAL_DIALOG);
  }, []);

  const confirm = useCallback((options) => new Promise((resolve) => {
    resolverRef.current?.(false);
    resolverRef.current = resolve;
    setDialog({ ...INITIAL_DIALOG, ...options, open: true });
  }), []);

  return (
    <ConfirmContext.Provider value={{ confirm }}>
      {children}
      <ConfirmDialog {...dialog} onConfirm={() => close(true)} onCancel={() => close(false)} />
    </ConfirmContext.Provider>
  );
}
