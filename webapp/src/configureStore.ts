import { createStore, compose } from 'redux'
import { reducer } from './store'

export default function configureStore(
    initialState?: any
) {
    const enhancer = (window as any)["__REDUX_DEVTOOLS_EXTENSION__"] ? (window as any)["__REDUX_DEVTOOLS_EXTENSION__"]()(createStore) : createStore;
    return enhancer(
        reducer,
        initialState,
        compose(),
    )
}
