
import { h } from 'preact';
import { useLocation } from 'react-router';

const Success = () => {
    const { state } = useLocation();
    const { msg } = state;

    return (
        <section class="Success">
            <div >
                Success !!
            </div>
            <div >
                {msg}
            </div>
        </section>
    );
};

export default Success;