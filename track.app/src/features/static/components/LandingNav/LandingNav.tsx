import { useNavigate } from "react-router-dom"
import styles from "./LandingNav.module.css"

const LandingNav = () => {
  const navigate = useNavigate()
  return (
    <div className={styles.section}>
      <h1>Brief<span>IQ</span></h1>
      <div className={styles.container}>
        <div className={styles.signin} onClick={() => navigate("/signin")}>
          <p>Sign in</p>
        </div>
        <div className={styles.signup} onClick={() => navigate("/signup")}>
          <p>Sign up</p>
        </div>
      </div>
    </div>
  )
}

export default LandingNav