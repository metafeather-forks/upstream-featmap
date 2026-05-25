import { configureStore } from '@reduxjs/toolkit'
import { reducer } from './store'

export default function createAppStore(initialState?: any) {
  return configureStore({
    reducer,
    preloadedState: initialState,
  })
}
