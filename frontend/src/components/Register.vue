<template>
    <div class="register__wrapper">
        <img src="../assets/pexels-cottonbro-5053739.jpg" alt="man holding phone">

        <div class="register">
            <h1>Create your account</h1>
            <form @submit.prevent="submitRegData" id="register__form">

                <div class="form-group">
                    <div class="form-input">
                        <label for="firstname">First name</label>
                        <input v-model="form.firstname" type="text" name="firstname" id="firstname" required>
                    </div>
                    <div class="form-input">
                        <label for="email">Email</label>
                        <input v-model="form.email" type="email" name="email" id="email" required>
                    </div>

                    <div class="form-input">
                        <label for="date">Date of Birth</label>
                        <input v-model="form.dateofbirth" type="date" name="date" id="date" required>
                    </div>

                    <div class="form-input">
                        <label for="aboutme">About me</label>
                        <textarea v-model="form.aboutme" id="aboutme" name="aboutme" rows="4" cols="50"></textarea>
                    </div>
                </div>

                <div class="form-group">
                    <div class="form-input">
                        <label for="lastname">Last name</label>
                        <input v-model="form.lastname" type="text" name="lastname" id="lastname" required>
                    </div>

                    <div class="form-input">
                        <label for="password">Password</label>
                        <input v-model="form.password" type="password" name="password" id="password" required>
                    </div>

                    <div class="form-input">
                        <label for="nickname">Nickname</label>
                        <input v-model="form.nickname" type="text" name="nickname" id="nickname">
                    </div>

                    <FileUpload v-model:file="form.avatar" labelName="Avatar"></FileUpload>
                </div>
            </form>

            <button class="btn" form="register__form" type="submit">Create account</button>

        </div>
    </div>

</template>

<script>
import FileUpload from './FileUpload.vue';

export default {
    name: "Register",
    data() {
        return {
            form: {
                email: "",
                password: "",
                firstname: "",
                lastname: "",
                dateofbirth: null,
                nickname: "",
                avatar: null,
                aboutme: "",
            },
        };
    },
    methods: {
        async submitRegData() {
            // for multipart form to work body should be of FormData type
            let formData = new FormData();
            formData.set("avatar", this.form.avatar);
            formData.set("email", this.form.email);
            formData.set("password", this.form.password);
            formData.set("firstname", this.form.firstname);
            formData.set("lastname", this.form.lastname);
            formData.set("dateofbirth", this.form.dateofbirth);
            formData.set("nickname", this.form.nickname);
            formData.set("aboutme", this.form.aboutme);
            
            await fetch("http://localhost:8081/register", {
                credentials: "include",
                method: "POST",
                body: formData
            })
                .then((res) => {
                    // console.log(res)
                    if (res.status === 409) {
                        this.$toast.open({
                            message: "Email already taken",
                            type: "error", //One of success, info, warning, error, default
                        });
                    }
                    else if (res.status === 400) {
                        this.$toast.open({
                            message: "Bad request",
                            type: "error", //One of success, info, warning, error, default
                        });
                    }
                    else {
                        this.$toast.open({
                            message: "Sucessfully registered!",
                            type: "success", //One of success, info, warning, error, default
                        });
                        this.$router.push("/");
                    }
                });
        },
    },
    components: { FileUpload }
}
</script>

<style>
.register__wrapper {
    display: flex;
    align-items: center;
    height: 80vh;
    min-height: 650px;
    max-height: 700px;
    width: 80vw;
    max-width: 1050px;
    min-width: 950px;
    margin: auto;
    background-color: var(--color-white);
    box-shadow: var(--container-shadow);
    border-radius: 20px;
    overflow: hidden;

}

.register {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    gap: 45px;
    padding: 0 50px;

}

.register__wrapper img {
    height: 100%;
}

.register form {
    display: flex;
    gap: 40px;
}



.register .form-group {
    gap: 25px;
    flex: 1;
}
</style>