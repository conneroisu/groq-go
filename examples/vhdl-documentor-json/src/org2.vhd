library IEEE;

use IEEE.std_logic_1164.all;

entity org2 is

  port(
    i_A : in  std_logic;
    i_B : in  std_logic;
    o_F : out std_logic
    );

end org2;

architecture dataflow of org2 is
begin

  o_F <= i_A or i_B;

end dataflow;
