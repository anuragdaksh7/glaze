import { useState } from "react";
import styles from "./SignUp.module.css"
import { useAppDispatch } from "../../../../app/hooks";
import { signUpUser } from "../../auth.thunk";
import { useNavigate } from "react-router-dom";
import { getGoogleRedirectURL } from "../../auth.api";

const SignUp = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const [email, setEmail] = useState("");
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");

  const handleSignUp = async () => {
    dispatch(signUpUser({ name, email, password }))
  }

  const handleGoogleLogin = async () => {
    const { redirect_url } = await getGoogleRedirectURL();
    window.location.href = redirect_url;
  }

  return (
    <div className={styles.section}>
      <main>
        <div className={styles.form}>
          <h1>Brief <span>IQ</span></h1>
          <div className={styles.fields}>
            <div className={styles.col}>
              <input type="text" placeholder="Enter your name..." value={name} onChange={(e) => setName(e.target.value)} />
              <input type="email" placeholder="Enter your email..." value={email} onChange={(e) => setEmail(e.target.value)} />
              <input type="password" placeholder="Enter your password..." value={password} onChange={(e) => setPassword(e.target.value)} />
            </div>
            <div className={styles.action}>
              <div className={styles.btns}>
                <div className={styles.btn} onClick={handleSignUp}>
                  <p>Get started!</p>
                </div>
                <div className={styles.googleBtn} onClick={handleGoogleLogin}>
                  <p>Continue with Google</p>
                </div>
              </div>
              <p>Already have account?<span onClick={() => navigate("/login")}> Login</span></p>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}

export default SignUp
