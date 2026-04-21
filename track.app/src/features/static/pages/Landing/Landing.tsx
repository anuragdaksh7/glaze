
import LandingNav from "../../components/LandingNav/LandingNav"
import styles from "./Landing.module.css"

const Landing = () => {
  return (
    <div className={styles.page}>
      <LandingNav />
      <div className={styles.landing}>
      </div>
    </div>
  )
}

export default Landing