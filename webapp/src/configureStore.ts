import { configureStore } from '@reduxjs/toolkit'
import { reducer, AppState } from './store'

export default function createAppStore(initialState?: Partial<AppState>) {
  return configureStore({
    reducer,
    preloadedState: initialState,
  })
}
