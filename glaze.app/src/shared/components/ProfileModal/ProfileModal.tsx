import { useState } from "react"
import { useAppSelector } from "../../../app/hooks"
import { COLORS } from "../../../constants/colors"
import { selectUser } from "../../../features/auth/auth.selector"
import { getRNGIndexBounded } from "../../../utils/rng-bounded"
import styles from "./ProfileModal.module.css"
import { useOutsideClickHandler } from "../../../hooks"
import { Logout, User } from "../../../icons"

const ProfileModal = () => {
  const APP_USER = useAppSelector(selectUser)

  const [modalOpen, setModalOpen] = useState(false)

  const modalHook = useOutsideClickHandler(() => {
    setModalOpen(false)
  })

  const initials = APP_USER?.name
  ?.trim()
  ?.split(/\s+/)
  ?.slice(0, 2)
  ?.map(w => w[0]?.toUpperCase())
  ?.join("") || "";
  return (
    <div className={styles.section} onClick={() => setModalOpen(!modalOpen)} style={{
      backgroundColor: COLORS[getRNGIndexBounded(0, COLORS.length - 1, initials)]
    }}>
      <p>{initials}</p>
      {
        modalOpen && <div ref={modalHook} className={styles.modal}>
          <div className={styles.row}>
            <User />
            <p>Profile</p>
          </div>
          <div className={`${styles.row} ${styles.logout}`}>
            <Logout />
            <p>Logout</p>
          </div>
        </div>
      }
    </div>
  )
}

export default ProfileModal