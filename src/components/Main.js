import React from 'react';
import { Route, Switch, Redirect } from 'react-router';

import Login from './Login';
import Register from './Register';
import Home from './Home';

function Main(props) {

    const {isLoggedIn, hanldeLoggedIn } = props;

    const showLogin = () => {
        return isLoggedIn ?
                (<Redirect to='/home' />) :
                (<Login hanldeLoggedIn={hanldeLoggedIn} />);
    };

    const showHome = () => {
        return isLoggedIn ? <Home /> : <Redirect to='/login' />
    };
    return (
        <div className="main">
            <Switch>
                <Route path='/' exact render={showLogin} />
                <Route path='/login' render={showLogin}></Route>
                <Route path='/register' component={Register}></Route>
                <Route path='/home' render={showHome}></Route>
            </Switch>
        </div>
    );
}

export default Main;