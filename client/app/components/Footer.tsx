import * as React from 'react';

const BFooter = require('re-bulma/lib/layout/footer').default;

interface BFooterProps extends React.HTMLAttributes {
}

const defaultStyle = {
    color: '#fff',
    backgroundColor: '#3572b0'
};

export function Footer(props: BFooterProps) {
    return <BFooter {...props} >{props.children}</BFooter>;
}
