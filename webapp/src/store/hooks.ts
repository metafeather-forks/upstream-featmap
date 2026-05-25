import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux'
import type { AppState, AllActions } from '../store'
import type { Dispatch } from 'redux'

export const useAppDispatch = () => useDispatch<Dispatch<AllActions>>()
export const useAppSelector: TypedUseSelectorHook<AppState> = useSelector
