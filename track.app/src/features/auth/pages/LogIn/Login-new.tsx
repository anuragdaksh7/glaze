import { useState } from "react";
import styles from "./Login-new.module.css"
import { useAppDispatch } from "../../../../app/hooks";
import { signInUser } from "../../auth.thunk";
import { useNavigate } from "react-router-dom";
import { getGoogleRedirectURL } from "../../auth.api";

const Login = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = async () => {
    dispatch(signInUser({ email, password }))
  }

  const handleGoogleLogin = async () => {
    const { redirect_url: redirectURL } = await getGoogleRedirectURL();
    console.log("Redirect URL:", redirectURL);
    window.location.href = redirectURL;
  }

  return (
    <div className={styles.section}>
      <main>
        <div className={styles.form}>
          <h1>Brief <span>IQ</span></h1>
          <div className={styles.fields}>
            <div className={styles.col}>
              <input type="email" placeholder="Enter your email..." value={email} onChange={(e) => setEmail(e.target.value)} />
              <input type="password" placeholder="Enter your password..." value={password} onChange={(e) => setPassword(e.target.value)} />
            </div>
            <div className={styles.action}>
              <div className={styles.btns}>
                <div className={styles.btn} onClick={handleLogin}>
                  <p>Login</p>
                </div>
                <div className={styles.googleBtn} onClick={handleGoogleLogin}>
                  <p>Continue with Google</p>
                </div>
              </div>
              <p>Need to create an account? <span onClick={() => navigate("/signup")}>Sign up</span></p>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}

export default Login
