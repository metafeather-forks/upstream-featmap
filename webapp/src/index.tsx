import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import { Provider } from 'react-redux'

import configureStore from './configureStore'
import { BrowserRouter as Router } from 'react-router-dom'

const store = configureStore()

ReactDOM.render(
    <Provider store={store}>
        <Router >
            <App application={{ mode: "", workspaces: [], memberships: [], messages: [] }} />
        </Router>
    </Provider>
    , document.getElementById('root'));


