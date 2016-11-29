import * as React from 'react';

import { Color } from './Modifiers';

const BAddons = require('re-bulma/lib/elements/addons').default;

interface AddonsProps extends React.HTMLAttributes {
    color?: Color;
}

export class Addons extends React.PureComponent<AddonsProps, void> {
    render() {
        return <BAddons {...this.props} />;
    }
}
