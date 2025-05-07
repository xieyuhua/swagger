<?php
namespace UsingRefs;

/**
 * @SWG\Definition()
 */
class Product {

    /**
     * The unique identifier of a product in our catalog.
     *
     * @var integer
     * @SWG\Property(format="int64")
     */
    public $id;

    /**
     * @SWG\Property(ref="#/definitions/product_status")
     */
    public $status;
}



/**
 * @SWG\Path(
 *   path="/products/{product_id}",
 *   @SWG\Parameter(ref="#/parameters/product_id_in_path_required")
 * )
 */

class ProductController {

    /**
     * @SWG\Get(
     *   tags={"Products"},
     *   path="/products/{product_id}",
     *   @SWG\Response(response="default", ref="#/responses/product")
     * )
     */
    public function getProduct($id) {

    }

    /**
     * @SWG\Post(
     *   tags={"Products"},
     *   path="/products/{product_id}",
     *   @SWG\Parameter(ref="#/parameters/product_in_body"),
     *   @SWG\Response(response="default", ref="#/responses/product")
     * )
     */
    public function updateProduct($id) {

    }

    /**
     * @SWG\Post(
     *   tags={"Products"},
     *   path="/products",
     *   @SWG\Parameter(ref="#/parameters/product_in_body"),
     *   @SWG\Response(response="default", ref="#/responses/product")
     * )
     */
    public function addProduct($id) {

    }

}
