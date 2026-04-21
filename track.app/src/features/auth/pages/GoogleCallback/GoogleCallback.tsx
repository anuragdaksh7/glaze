import { useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useAppDispatch } from "../../../../app/hooks";
import { verifyGoogleCode } from "../../auth.thunk";

const GoogleCallback = () => {
  const location = useLocation();
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  useEffect(() => {
    const params = new URLSearchParams(location.search);
    const code = params.get("code");

    if (code) {
      dispatch(verifyGoogleCode(code))
        .unwrap()
        .then(() => {
          navigate("/audio");
        })
        .catch((err) => {
          console.error("Google OAuth verification failed:", err);
          navigate("/login");
        });
    } else {
      navigate("/login");
    }
  }, [location, dispatch, navigate]);

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
      <p>Verifying Google account...</p>
    </div>
  );
};

export default GoogleCallback;
