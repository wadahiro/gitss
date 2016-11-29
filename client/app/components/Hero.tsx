import * as React from 'react';

import { Size } from './Modifiers';

const BHero = require('re-bulma/lib/layout/hero').default;
const BHeroBody = require('re-bulma/lib/layout/hero-body').default;
const BHeroFoot = require('re-bulma/lib/layout/hero-foot').default;
const BHeroHead = require('re-bulma/lib/layout/hero-head').default;


interface HeroProps extends React.HTMLAttributes {
    color?: 'isPrimary' | 'isDanger';
}

const heroStyle = {
    backgroundColor: '#205081'
}

export class Hero extends React.PureComponent<HeroProps, void> {
    render() {
        return <BHero {...this.props} style={{ ...heroStyle, ...this.props.style }} />;
    }
}

export class HeroHead extends React.PureComponent<HeroProps, void> {
    render() {
        return <BHeroHead {...this.props} />;
    }
}

export class HeroBody extends React.PureComponent<HeroProps, void> {
    render() {
        return <BHeroBody {...this.props} />;
    }
}

export class HeroFoot extends React.PureComponent<HeroProps, void> {
    render() {
        return <BHeroFoot {...this.props} />;
    }
}