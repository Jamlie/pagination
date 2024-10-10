// @ts-check

const router = require("express").Router();
const { insert, show, paginate } = require("../controller/pagination");

router.post("/insert", insert);
router.get("/show", show);
router.post("/paginate", paginate);

module.exports = router;
