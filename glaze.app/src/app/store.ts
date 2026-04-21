import { configureStore } from '@reduxjs/toolkit';
// import tasksReducer from '../features/auth/auth.slice';
import userReducer from '../features/auth/auth.slice';

export const store = configureStore({
  reducer: {
    user: userReducer,
  },
  devTools: import.meta.env.DEV,
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
