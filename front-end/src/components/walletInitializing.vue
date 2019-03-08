<template>
    <div class="preLoading">
        <div>
            Please waiting for about 5 Seconds
        </div>
        <div class="loading">
            <div>
                Initializing Wallet
            </div>
            <div class="pacman">
                <pacman> </pacman>
            </div>
        </div>
    </div>
</template>

<script>
    // @ is an alias to /src
    require("@/assets/css/style.css");
    export default {
        name: "walletLoading",
        data: function () {
            return {
            }
        },
    };
</script>
